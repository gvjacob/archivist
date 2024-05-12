package storage

import (
	"database/sql"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
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

	const checkUserQuery = `
    SELECT COUNT(*) FROM users;
  `

	const insertUserQuery = `
    INSERT INTO users (access_token, refresh_token, last_archived) VALUES (?, ?, ?);
  `

	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(initQuery); err != nil {
		return nil, err
	}

	var count int
	if err := db.QueryRow(checkUserQuery).Scan(&count); err != nil {
		return nil, err
	}

	if count == 0 {
		log.Warn("No user found, creating a new one")

		_, err := db.Exec(
			insertUserQuery,
			os.Getenv("SPOTIFY_ACCESS_TOKEN"),
			os.Getenv("SPOTIFY_REFRESH_TOKEN"),
			time.Now().UTC().Format(time.RFC3339),
		)

		if err != nil {
			return nil, err
		}
	}

	return &Database{db: db}, nil
}
