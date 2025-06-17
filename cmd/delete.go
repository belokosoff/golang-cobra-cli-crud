package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/belokosoff/golang-cobra-cli-crud/internal/repository"
	"github.com/belokosoff/golang-cobra-cli-crud/pkg/db"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a book by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.InitDB()
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()

		id, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid ID format: %v", err)
		}

		repo := repository.NewBookRepository(db)
		err = repo.DeleteBook(id)
		if err != nil {
			log.Fatalf("Failed to delete book: %v", err)
		}

		fmt.Printf("Book with ID %d deleted successfully\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
