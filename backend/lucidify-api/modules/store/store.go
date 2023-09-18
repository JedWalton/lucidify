package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore(postgresqlURL string) (*Store, error) {
	if postgresqlURL == "" {
		return nil, fmt.Errorf("POSTGRESQL_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", postgresqlURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}
