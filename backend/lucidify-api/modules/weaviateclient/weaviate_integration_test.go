// //go:build integration
// // +build integration
package weaviateclient

import "testing"

func TestWeaviateClient(t *testing.T) {
	client, err := NewWeaviateClient()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	// Test uploading a document
	err = client.UploadDocument("testuser", "testdoc", "test content")
	if err != nil {
		t.Errorf("failed to upload document: %v", err)
	}

	// Test updating a document
	err = client.UpdateDocument("testuser", "testdoc", "updated test content")
	if err != nil {
		t.Errorf("failed to update document: %v", err)
	}

	// Test deleting a document
	err = client.DeleteDocument("testuser", "testdoc")
	if err != nil {
		t.Errorf("failed to delete document: %v", err)
	}
}
