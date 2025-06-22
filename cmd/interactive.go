package cmd

import (
	"database/sql"
	"log"

	"github.com/belokosoff/golang-cobra-cli-crud/tui"
	"github.com/spf13/cobra"
)

var interactiveCommand = &cobra.Command{
	Use:   "interactive",
	Short: "Run TUI mode of application",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "./books.db")
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		tui.Start(db)
	},
}

func init() {
	rootCmd.AddCommand(interactiveCommand)
}
