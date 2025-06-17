package repository

import (
	"database/sql"
	"fmt"

	"github.com/belokosoff/golang-cobra-cli-crud/internal/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) GetAllBooks() ([]models.Book, error) {
	rows, err := r.db.Query("SELECT id, title, author, published_year, status FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.Status)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) GetFilteredBooks(filter string) ([]models.Book, error) {
	query := "SELECT id, title, author, published_year, status FROM books WHERE status = ?"
	rows, err := r.db.Query(query, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.PublishedYear, &b.Status)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) AddBook(book models.Book) error {
	query := `INSERT INTO books (title, author, published_year, status) VALUES (?, ?, ?, ?)`
	_, err := r.db.Exec(query, book.Title, book.Author, book.PublishedYear, book.Status)
	return err
}

func (r *BookRepository) UpdateStatusBook(id int) error {
	query := `UPDATE books SET status = 'read' WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("book with ID %d not found", id)
	}
	return nil
}

func (r *BookRepository) DeleteBook(id int) error {
	query := `DELETE FROM books WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("book with ID %d not found", id)
	}

	return nil
}
