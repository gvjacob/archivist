package storage

import (
	"database/sql"
	"sync"

	_ "modernc.org/sqlite"
)

type Database struct {
	mu sync.Mutex
	*sql.DB
}

func NewDatabase(file string) (*Database, error) {
	db, err := sql.Open("sqlite", file)

	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}
