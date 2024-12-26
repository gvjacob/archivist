package main

import (
	"archivist/utils"
	"database/sql"
	"os"

	"github.com/charmbracelet/log"
	_ "modernc.org/sqlite"
)

func main() {
	utils.LoadDotEnv()

	log.Info("Creating tables...")

	if err := createTables(); err != nil {
		log.Fatal(err)
	}
}

func createTables() error {
	const createTablesQuery = `
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      access_token TEXT NOT NULL,
      refresh_token TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS archived_tracks (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER NOT NULL,
      track_id TEXT NOT NULL,
      playlist_id TEXT NULL,
      created_at INTEGER NOT NULL,
      FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
    );
  `

	db, err := sql.Open("sqlite", os.Getenv("SQLITE_FILE_PATH"))

	if err != nil {
		return err
	}

	if _, err := db.Exec(createTablesQuery); err != nil {
		return err
	}

	return nil
}
