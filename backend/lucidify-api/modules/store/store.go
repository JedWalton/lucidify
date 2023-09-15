package store

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func NewStore() *Store {
	connectionString := os.Getenv("POSTGRESQL_URL")
	if connectionString == "" {
		log.Fatal("POSTGRESQL_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return &Store{db: db}
}
