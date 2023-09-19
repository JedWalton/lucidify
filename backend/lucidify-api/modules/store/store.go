package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
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

func SetupTestStore() (*Store, error) {
	// Get the path to the current file (store.go)
	_, filename, _, _ := runtime.Caller(0)

	// Get the directory containing store.go
	dir := filepath.Dir(filename)

	// Construct the path to the .env file relative to store.go
	envPath := filepath.Join(dir, "../../../../.env")

	// Load environment variables from the .env file
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Get the PostgreSQL URL from the environment variables
	postgresqlURL := os.Getenv("POSTGRESQL_URL")

	// Create a new store instance using the NewStore function
	storeInstance, err := NewStore(postgresqlURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %v", err)
	}

	return storeInstance, nil
}
