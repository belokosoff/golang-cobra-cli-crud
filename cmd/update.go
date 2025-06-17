package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/belokosoff/golang-cobra-cli-crud/internal/repository"
	"github.com/belokosoff/golang-cobra-cli-crud/pkg/db"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update status a book by ID",
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
		err = repo.UpdateStatusBook(id)
		if err != nil {
			log.Fatalf("Failed to update book: %v", err)
		}

		fmt.Printf("Status book with ID %d update successfully\n", id)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
