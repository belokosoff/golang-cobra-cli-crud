package cmd

import (
	"fmt"
	"log"

	"github.com/belokosoff/golang-cobra-cli-crud/internal/repository"
	"github.com/belokosoff/golang-cobra-cli-crud/pkg/db"
	"github.com/spf13/cobra"
)

var findByIdStatusCmd = &cobra.Command{
	Use:   "find-by-status",
	Short: "Find books by status (read/unread)",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.InitDB()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		repo := repository.NewBookRepository(db)

		status, _ := cmd.Flags().GetString("status")
		books, err := repo.GetFilteredBooks(status)
		if err != nil {
			log.Fatalf("Failed to find books: %v", err)
		}

		if len(books) == 0 {
			fmt.Println("No books found with status:", status)
			return
		}

		fmt.Printf("Books with status '%s':\n", status)
		for _, book := range books {
			fmt.Printf("- ID: %d, Title: %s, Author: %s, Year: %d\n",
				book.ID, book.Title, book.Author, book.PublishedYear)
		}
	},
}

func init() {
	rootCmd.AddCommand(findByIdStatusCmd)
	findByIdStatusCmd.Flags().StringP("status", "s", "", "Filter by status (read/unread)")
	findByIdStatusCmd.MarkFlagRequired("status")
}
