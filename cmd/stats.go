package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show book statistics",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sql.Open("sqlite3", "./books.db")
		if err != nil {
			log.Fatalf("Failed to open database: %v", err)
		}
		defer db.Close()

		byYear, _ := cmd.Flags().GetBool("by-year")
		byAuthor, _ := cmd.Flags().GetBool("by-author")
		byStatus, _ := cmd.Flags().GetBool("by-status")

		if !byYear && !byAuthor && !byStatus {
			showBasicStats(db)
			return
		}

		if byYear {
			showYearStats(db)
		}
		if byAuthor {
			showAuthorStats(db)
		}
		if byStatus {
			showStatusStats(db)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)

	statsCmd.Flags().BoolP("by-year", "y", false, "Show statistics by publication year")
	statsCmd.Flags().BoolP("by-author", "a", false, "Show statistics by author")
	statsCmd.Flags().BoolP("by-status", "s", false, "Show read/unread statistics")
}

func showBasicStats(db *sql.DB) {
	var total, read int
	err := db.QueryRow("SELECT COUNT(*) FROM books").Scan(&total)
	if err != nil {
		log.Fatal(err)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM books WHERE status = 'read'").Scan(&read)
	if err != nil {
		log.Fatal(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nSTATISTIC\tVALUE\t")
	fmt.Fprintln(w, "---------\t-----\t")
	fmt.Fprintf(w, "Total books\t%d\t\n", total)
	fmt.Fprintf(w, "Read\t%d (%.0f%%)\t\n", read, float64(read)/float64(total)*100)
	fmt.Fprintf(w, "Unread\t%d (%.0f%%)\t\n", total-read, float64(total-read)/float64(total)*100)
	w.Flush()
}

func showYearStats(db *sql.DB) {
	rows, err := db.Query(`
		SELECT published_year, COUNT(*) as count 
		FROM books 
		GROUP BY published_year 
		ORDER BY published_year DESC`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nYEAR\tCOUNT\t")
	fmt.Fprintln(w, "----\t-----\t")

	var year, count int
	for rows.Next() {
		err := rows.Scan(&year, &count)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%d\t%d\t\n", year, count)
	}
	w.Flush()
}

func showAuthorStats(db *sql.DB) {
	rows, err := db.Query(`
		SELECT author, COUNT(*) as count 
		FROM books 
		GROUP BY author 
		ORDER BY count DESC`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nAUTHOR\tCOUNT\t")
	fmt.Fprintln(w, "------\t-----\t")

	var author string
	var count int
	for rows.Next() {
		err := rows.Scan(&author, &count)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\t%d\t\n", author, count)
	}
	w.Flush()
}

func showStatusStats(db *sql.DB) {
	rows, err := db.Query(`
		SELECT status, COUNT(*) as count 
		FROM books 
		GROUP BY status`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nSTATUS\tCOUNT\t")
	fmt.Fprintln(w, "------\t-----\t")

	var status string
	var count int
	for rows.Next() {
		err := rows.Scan(&status, &count)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\t%d\t\n", status, count)
	}
	w.Flush()
}
