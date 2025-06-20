package cmd

import (
	"fmt"
	"log"

	"github.com/belokosoff/golang-cobra-cli-crud/internal/models"
	"github.com/belokosoff/golang-cobra-cli-crud/internal/repository"
	"github.com/belokosoff/golang-cobra-cli-crud/pkg/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new book",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.InitDB()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		repo := repository.NewBookRepository(db)

		title, _ := cmd.Flags().GetString("title")
		author, _ := cmd.Flags().GetString("author")
		year, _ := cmd.Flags().GetInt("year")
		status, _ := cmd.Flags().GetString("status")

		book := models.Book{
			Title:         title,
			Author:        author,
			PublishedYear: year,
			Status:        status,
		}

		if err := repo.AddBook(book); err != nil {
			log.Fatalf("Failed to add book: %v", err)
		}
		fmt.Println("Book added successfully!")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("title", "t", "", "Book title")
	addCmd.Flags().StringP("author", "a", "", "Book author")
	addCmd.Flags().StringP("status", "s", "", "Book status (read/unread)")
	addCmd.Flags().IntP("year", "y", 0, "Published year")
}
