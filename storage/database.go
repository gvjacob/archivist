package storage

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	mu sync.Mutex
	db *sql.DB
}

func NewDatabase(file string) (*Database, error) {
	const initQuery = `
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      access_token TEXT NOT NULL,
      refresh_token TEXT NOT NULL,
      last_archived TEXT
    );
  `

	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(initQuery); err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}
