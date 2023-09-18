//go:build integration
// +build integration

package store

import (
	"lucidify-api/modules/testutils"
	"testing"

	_ "github.com/lib/pq"
)

func TestIntegration_UploadDocument(t *testing.T) {
	db := testutils.SetupDB()
	defer db.Close()

	store := &Store{db: db}

	// Test
	document_name := "Test Document"
	content := "This is a test document content."

	err := store.UploadDocument(document_name, content)
	if err != nil {
		t.Fatalf("Failed to upload document: %v", err)
	}

	// Verify
	var query_res_name, query_res_content string
	err = store.db.QueryRow("SELECT document_name, content FROM documents WHERE document_name = $1", document_name).Scan(&query_res_name, &query_res_content)
	if err != nil {
		t.Fatalf("Failed to retrieve document: %v", err)
	}

	if query_res_name != document_name || query_res_content != content {
		t.Fatalf("Document mismatch. Expected: (%s, %s). Got: (%s, %s)", document_name, content, query_res_name, query_res_content)
	}

	// Cleanup
	_, err = store.db.Exec("DELETE FROM documents WHERE document_name = $1", document_name)
	if err != nil {
		t.Fatalf("Failed to clean up test document: %v", err)
	}
}

func TestGetDocument(t *testing.T) {
	db := testutils.SetupDB()
	defer db.Close()

	store := &Store{db: db}

	// Test Data
	document_name := "Test Document for Retrieval"
	content := "This is content for the retrieval test."

	// Insert test document
	_, err := store.db.Exec(`INSERT INTO documents (document_name, content) VALUES ($1, $2)`, document_name, content)
	if err != nil {
		t.Fatalf("Failed to insert test document: %v", err)
	}

	// Test
	retrievedContent, err := store.GetDocument(document_name)
	if err != nil {
		t.Fatalf("Failed to retrieve document: %v", err)
	}

	// Verify
	if retrievedContent != content {
		t.Fatalf("Document content mismatch. Expected: %s. Got: %s", content, retrievedContent)
	}

	// Cleanup
	_, err = store.db.Exec("DELETE FROM documents WHERE document_name = $1", document_name)
	if err != nil {
		t.Fatalf("Failed to clean up test document: %v", err)
	}
}

func TestDeleteDocument(t *testing.T) {
	db := testutils.SetupDB()
	defer db.Close()

	store := &Store{db: db}

	// Test Data
	document_name := "Test Document for Deletion"
	content := "This is content for the deletion test."

	// Insert test document
	_, err := store.db.Exec(`INSERT INTO documents (document_name, content) VALUES ($1, $2)`, document_name, content)
	if err != nil {
		t.Fatalf("Failed to insert test document: %v", err)
	}

	// Test
	err = store.DeleteDocument(document_name)
	if err != nil {
		t.Fatalf("Failed to delete document: %v", err)
	}

	// Verify
	var count int
	err = store.db.QueryRow(`SELECT COUNT(*) FROM documents WHERE document_name = $1`, document_name).Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query document count: %v", err)
	}
	if count != 0 {
		t.Fatalf("Document was not deleted. Expected count: 0. Got: %d", count)
	}
}

func TestUpdateDocument(t *testing.T) {
	db := testutils.SetupDB()
	defer db.Close()

	store := &Store{db: db}

	// Test Data
	document_name := "Test Document for Update"
	original_content := "This is the original content."
	updated_content := "This is the updated content."

	// Insert test document
	_, err := store.db.Exec(`INSERT INTO documents (document_name, content) VALUES ($1, $2)`, document_name, original_content)
	if err != nil {
		t.Fatalf("Failed to insert test document: %v", err)
	}

	// Test
	err = store.UpdateDocument(document_name, updated_content)
	if err != nil {
		t.Fatalf("Failed to update document: %v", err)
	}

	// Verify
	var retrievedContent string
	err = store.db.QueryRow(`SELECT content FROM documents WHERE document_name = $1`, document_name).Scan(&retrievedContent)
	if err != nil {
		t.Fatalf("Failed to retrieve updated document content: %v", err)
	}
	if retrievedContent != updated_content {
		t.Fatalf("Document content mismatch. Expected: %s. Got: %s", updated_content, retrievedContent)
	}

	// Cleanup
	_, err = store.db.Exec("DELETE FROM documents WHERE document_name = $1", document_name)
	if err != nil {
		t.Fatalf("Failed to clean up test document: %v", err)
	}
}
