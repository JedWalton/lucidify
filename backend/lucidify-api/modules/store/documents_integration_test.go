// //go:build integration
// // +build integration
package store

import (
	"lucidify-api/modules/config"
	"testing"

	_ "github.com/lib/pq"
)

func (s *Store) createTestUser(db *Store) (string, error) {
	userID := "testuuid1237fyuiaroi"
	const query = `INSERT INTO users (user_id, username, email) VALUES ($1, 'testuser', 'test@example.com') RETURNING user_id`
	if err := s.db.QueryRow(query, userID).Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func (s *Store) deleteTestUser(db *Store, userID string) error { // Changed parameter type
	_, err := s.db.Exec(`DELETE FROM users WHERE user_id = $1`, userID)
	return err
}

func (s *Store) TestIntegration_UploadDocument(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewStore(PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Create test user
	userID, err := s.createTestUser(store)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test
	document_name := "Test Document"
	content := "This is a test document content."

	err = store.UploadDocument(userID, document_name, content)
	if err != nil {
		t.Fatalf("Failed to upload document: %v", err)
	}

	// Verify
	var query_res_name, query_res_content string
	err = store.db.QueryRow("SELECT document_name, content FROM documents WHERE user_id = $1 AND document_name = $2", userID, document_name).Scan(&query_res_name, &query_res_content)
	if err != nil {
		t.Fatalf("Failed to retrieve document: %v", err)
	}

	if query_res_name != document_name || query_res_content != content {
		t.Fatalf("Document mismatch. Expected: (%s, %s). Got: (%s, %s)", document_name, content, query_res_name, query_res_content)
	}

	// Cleanup
	_, err = store.db.Exec("DELETE FROM documents WHERE user_id = $1 AND document_name = $2", userID, document_name)
	if err != nil {
		t.Fatalf("Failed to clean up test document: %v", err)
	}

	// Delete test user
	err = store.deleteTestUser(store, userID)
	if err != nil {
		t.Fatalf("Failed to delete test user: %v", err)
	}
}

func (s *Store) TestGetDocument(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewStore(PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Create test user
	userID, err := store.createTestUser(store)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test Data
	document_name := "Test Document for Retrieval"
	content := "This is content for the retrieval test."

	// Insert test document
	_, err = store.db.Exec(`INSERT INTO documents (user_id, document_name, content) VALUES ($1, $2, $3)`, userID, document_name, content)
	if err != nil {
		t.Fatalf("Failed to insert test document: %v", err)
	}

	// Test
	retrievedContent, err := store.GetDocument(userID, document_name)
	if err != nil {
		t.Fatalf("Failed to retrieve document: %v", err)
	}

	// Verify
	if retrievedContent != content {
		t.Fatalf("Document content mismatch. Expected: %s. Got: %s", content, retrievedContent)
	}

	// Cleanup
	_, err = store.db.Exec("DELETE FROM documents WHERE user_id = $1 AND document_name = $2", userID, document_name)
	if err != nil {
		t.Fatalf("Failed to clean up test document: %v", err)
	}

	// Delete test user
	err = store.deleteTestUser(store, userID)
	if err != nil {
		t.Fatalf("Failed to delete test user: %v", err)
	}
}

func TestDeleteDocument(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewStore(PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Create test user
	userID, err := store.createTestUser(store)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test Data
	document_name := "Test Document for Deletion"
	content := "This is content for the deletion test."

	// Insert test document
	_, err = store.db.Exec(`INSERT INTO documents (user_id, document_name, content) VALUES ($1, $2, $3)`, userID, document_name, content)
	if err != nil {
		t.Fatalf("Failed to insert test document: %v", err)
	}

	// Test
	err = store.DeleteDocument(userID, document_name)
	if err != nil {
		t.Fatalf("Failed to delete document: %v", err)
	}

	// Verify
	var count int
	err = store.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE user_id = $1 AND document_name = $2`, userID, document_name).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query document count: %v", err)
	}
	if count != 0 {
		t.Fatalf("Document was not deleted. Expected count: 0. Got: %d", count)
	}

	// Delete test user
	err = store.deleteTestUser(store, userID)
	if err != nil {
		t.Fatalf("Failed to delete test user: %v", err)
	}
}

func TestUpdateDocument(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewStore(PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Create test user
	userID, err := store.createTestUser(store)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Test Data
	document_name := "Test Document for Update"
	original_content := "This is the original content."
	updated_content := "This is the updated content."

	// Insert test document
	_, err = store.db.Exec(`INSERT INTO documents (user_id, document_name, content) VALUES ($1, $2, $3)`, userID, document_name, original_content)
	if err != nil {
		t.Fatalf("Failed to insert test document: %v", err)
	}

	// Test
	err = store.UpdateDocument(userID, document_name, updated_content)
	if err != nil {
		t.Fatalf("Failed to update document: %v", err)
	}

	// Verify
	var retrievedContent string
	err = store.db.QueryRow(`SELECT content FROM documents WHERE user_id = $1 AND document_name = $2`, userID, document_name).Scan(&retrievedContent)
	if err != nil {
		t.Fatalf("Failed to retrieve updated document content: %v", err)
	}
	if retrievedContent != updated_content {
		t.Fatalf("Document content mismatch. Expected: %s. Got: %s", updated_content, retrievedContent)
	}

	// Cleanup
	_, err = store.db.Exec("DELETE FROM documents WHERE user_id = $1 AND document_name = $2", userID, document_name)
	if err != nil {
		t.Fatalf("Failed to clean up test document: %v", err)
	}

	// Delete test user
	err = store.deleteTestUser(store, userID)
	if err != nil {
		t.Fatalf("Failed to delete test user: %v", err)
	}
}
