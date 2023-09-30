//go:build integration
// +build integration

package postgresqlclient

import (
	"testing"
)

func TestIntegrationNewStore(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create postgresqlclient: %v", err)
	}

	if store.db == nil {
		t.Fatal("Expected db to be initialized, but it was nil")
	}

	// Teardown
	store.db.Close()
}
