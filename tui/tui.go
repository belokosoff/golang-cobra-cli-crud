package tui

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
)

type model struct {
	db          *sql.DB
	books       []Book
	cursor      int
	view        string
	title       string
	author      string
	year        string
	status      string
	activeField int // 0: title, 1: author, 2: year, 3: status
}

type Book struct {
	ID            int
	Title         string
	Author        string
	PublishedYear int
	Status        string
}

func initialModel(db *sql.DB) model {
	books := fetchBooks(db)
	return model{
		db:     db,
		books:  books,
		view:   "list",
		status: "unread",
	}
}

func fetchBooks(db *sql.DB) []Book {
	rows, err := db.Query("SELECT id, title, author, published_year, status FROM books")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.Status)
		if err != nil {
			log.Fatal(err)
		}
		books = append(books, b)
	}
	return books
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.view == "add" && msg != nil {
		if errMsg, ok := msg.(error); ok {
			log.Println("Error:", errMsg)
			return m, nil
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Обработка команд, которые работают в любом режиме
		switch msg.String() {
		case "ctrl+c", "esc":
			if m.view == "add" {
				m.view = "list"
				m.books = fetchBooks(m.db)
			} else {
				return m, tea.Quit
			}
			return m, nil
		}

		// Обработка команд для конкретных режимов
		switch m.view {
		case "list":
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.books)-1 {
					m.cursor++
				}
			case "a":
				m.view = "add"
				m.activeField = 0
				m.title = ""
				m.author = ""
				m.year = ""
				m.status = "unread"
			case "s":
				m.view = "stats"
			case "enter":
				// В режиме списка enter не делает ничего
			case "d":
				if len(m.books) > 0 {
					_, err := m.db.Exec("DELETE FROM books WHERE id = ?", m.books[m.cursor].ID)
					if err != nil {
						log.Println("Error deleting book:", err)
					}
					m.books = fetchBooks(m.db)
					if m.cursor >= len(m.books) {
						m.cursor = len(m.books) - 1
					}
				}
			case "t":
				if len(m.books) > 0 {
					newStatus := "read"
					if m.books[m.cursor].Status == "read" {
						newStatus = "unread"
					}
					_, err := m.db.Exec(
						"UPDATE books SET status = ? WHERE id = ?",
						newStatus, m.books[m.cursor].ID,
					)
					if err != nil {
						log.Println("Error updating status:", err)
					}
					m.books = fetchBooks(m.db)
				}
			}

		case "add":
			switch msg.String() {
			case "tab":
				m.activeField = (m.activeField + 1) % 4
			case "shift+tab":
				m.activeField = (m.activeField - 1 + 4) % 4
			case "enter":
				if strings.TrimSpace(m.title) == "" {
					log.Println("Title cannot be empty")
					return m, nil
				}
				if strings.TrimSpace(m.author) == "" {
					log.Println("Author cannot be empty")
					return m, nil
				}
				if strings.TrimSpace(m.year) == "" {
					log.Println("Year cannot be empty")
					return m, nil
				}

				year, err := strconv.Atoi(m.year)
				if err != nil {
					log.Println("Invalid year format")
					return m, nil
				}

				_, err = m.db.Exec(
					"INSERT INTO books (title, author, published_year, status) VALUES (?, ?, ?, ?)",
					strings.TrimSpace(m.title),
					strings.TrimSpace(m.author),
					year,
					m.status,
				)
				if err != nil {
					log.Println("Error adding book:", err)
				}
				m.view = "list"
				m.books = fetchBooks(m.db)
				m.title = ""
				m.author = ""
				m.year = ""
			case " ":
				if m.activeField == 3 { // Только для поля статуса
					if m.status == "read" {
						m.status = "unread"
					} else {
						m.status = "read"
					}
				} else {
					// Для других полей пробел - обычный символ
					switch m.activeField {
					case 0:
						m.title += " "
					case 1:
						m.author += " "
					case 2:
						m.year += " "
					}
				}
			case "backspace":
				switch m.activeField {
				case 0:
					if len(m.title) > 0 {
						m.title = m.title[:len(m.title)-1]
					}
				case 1:
					if len(m.author) > 0 {
						m.author = m.author[:len(m.author)-1]
					}
				case 2:
					if len(m.year) > 0 {
						m.year = m.year[:len(m.year)-1]
					}
				}
			default:
				// Обработка обычного ввода текста
				if len(msg.String()) == 1 {
					switch m.activeField {
					case 0:
						m.title += msg.String()
					case 1:
						m.author += msg.String()
					case 2:
						if _, err := strconv.Atoi(msg.String()); err == nil {
							m.year += msg.String()
						}
					}
				}
			}

		case "stats":
			// В режиме статистики не обрабатываем специальные команды
		}
	}

	return m, nil
}

func (m model) View() string {
	var sb strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
	selectedStyle := lipgloss.NewStyle().Background(lipgloss.Color("62")).Foreground(lipgloss.Color("230"))
	normalStyle := lipgloss.NewStyle().PaddingLeft(2)
	readStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	unreadStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	activeFieldStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	//errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)

	switch m.view {
	case "list":
		sb.WriteString(titleStyle.Render("Your Book Collection\n"))
		for i, book := range m.books {
			status := readStyle.Render("✓ ")
			if book.Status == "unread" {
				status = unreadStyle.Render("✗ ")
			}

			if m.cursor == i {
				sb.WriteString(selectedStyle.Render(fmt.Sprintf(
					"%s %s by %s (%d)",
					status, book.Title, book.Author, book.PublishedYear,
				)))
			} else {
				sb.WriteString(normalStyle.Render(fmt.Sprintf(
					"%s %s by %s (%d)",
					status, book.Title, book.Author, book.PublishedYear,
				)))
			}
			sb.WriteString("\n")
		}

		sb.WriteString("\n" + helpStyle.Render(
			"↑/↓: Navigate • a: Add • d: Delete • t: Toggle status • s: Stats • q: Quit",
		))

	case "add":
		sb.WriteString(titleStyle.Render("Add New Book\n\n"))

		// Title field
		titleLabel := "Title:"
		if m.activeField == 0 {
			titleLabel = activeFieldStyle.Render(titleLabel)
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", titleLabel, m.title))

		// Author field
		authorLabel := "Author:"
		if m.activeField == 1 {
			authorLabel = activeFieldStyle.Render(authorLabel)
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", authorLabel, m.author))

		// Year field
		yearLabel := "Year:"
		if m.activeField == 2 {
			yearLabel = activeFieldStyle.Render(yearLabel)
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", yearLabel, m.year))

		// Status field
		statusLabel := "Status:"
		if m.activeField == 3 {
			statusLabel = activeFieldStyle.Render(statusLabel)
		}
		statusValue := m.status
		if m.activeField == 3 {
			statusValue = activeFieldStyle.Render(statusValue)
		}
		sb.WriteString(fmt.Sprintf("%s %s\n\n", statusLabel, statusValue))

		sb.WriteString(helpStyle.Render(
			"Tab/Shift+Tab: Move between fields • Space: Toggle status • Enter: Save • Esc: Cancel",
		))

	case "stats":
		var total, read int
		m.db.QueryRow("SELECT COUNT(*) FROM books").Scan(&total)
		m.db.QueryRow("SELECT COUNT(*) FROM books WHERE status = 'read'").Scan(&read)

		sb.WriteString(titleStyle.Render("Statistics\n\n"))
		sb.WriteString(fmt.Sprintf("Total books: %d\n", total))
		sb.WriteString(fmt.Sprintf("Read:       %d (%.0f%%)\n", read, float64(read)/float64(total)*100))
		sb.WriteString(fmt.Sprintf("Unread:     %d (%.0f%%)\n\n", total-read, float64(total-read)/float64(total)*100))
		sb.WriteString(helpStyle.Render("Esc: Back to list"))
	}

	return sb.String()
}

func Start(db *sql.DB) {
	p := tea.NewProgram(initialModel(db))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
