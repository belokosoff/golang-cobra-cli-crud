package cmd

import (
	"fmt"
	"log"

	"github.com/belokosoff/golang-cobra-cli-crud/internal/repository"
	"github.com/belokosoff/golang-cobra-cli-crud/pkg/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Output the list of book",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.InitDB()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		repo := repository.NewBookRepository(db)

		books, err := repo.GetAllBooks()
		if err != nil {
			log.Fatalf("Failed to find books: %v", err)
		}

		if len(books) == 0 {
			fmt.Println("No books found")
			return
		}

		for _, book := range books {
			fmt.Printf("- ID: %d, Title: %s, Author: %s, Year: %d, Status: %s\n",
				book.ID, book.Title, book.Author, book.PublishedYear, book.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
