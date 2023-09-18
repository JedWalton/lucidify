//go:build integration
// +build integration

package store

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestIntegrationNewStore(t *testing.T) {
	// Load environment variables from .env file
	godotenv.Load("../../../../.env")

	// Get the POSTGRESQL_URL environment variable
	postgresqlURL := os.Getenv("POSTGRESQL_URL")

	store, err := NewStore(postgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}

	if store.db == nil {
		t.Fatal("Expected db to be initialized, but it was nil")
	}

	// Teardown
	store.db.Close()
}
