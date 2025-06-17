package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		author TEXT NOT NULL,
		published_year INTEGER,
		status TEXT DEFAULT 'unread'
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	//log.Println("Connected to SQLite database")
	return db, nil
}
