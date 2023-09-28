// //go:build integration
// // +build integration
package weaviateclient

import (
	"testing"

	"github.com/google/uuid"
)

func TestWeaviateClient(t *testing.T) {
	weaviateClient, err := NewWeaviateClient()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	// Test uploading a document
	documentID := uuid.New().String()
	err = weaviateClient.UploadDocument(documentID, "testuser", "testdoc", "test content")
	if err != nil {
		t.Errorf("failed to upload document: %v", err)
	}

	document, err := weaviateClient.GetDocument(documentID)
	if err != nil {
		t.Errorf("failed to get document: %v", err)
	}
	t.Logf("document: %+v", document)
	if document.UserID != "testuser" {
		t.Errorf("document owner is incorrect: %v", document.UserID)
	}
	if document.DocumentName != "testdoc" {
		t.Errorf("document name is incorrect: %v", document.DocumentName)
	}
	if document.Content != "test content" {
		t.Errorf("document content is incorrect: %v", document.Content)
	}

	// Test updating a document
	err = weaviateClient.UpdateDocumentContent(documentID, "updated test content")
	if err != nil {
		t.Errorf("failed to update document: %v", err)
	}
	document, err = weaviateClient.GetDocument(documentID)
	if document.Content != "updated test content" {
		t.Errorf("document content is incorrect: %v", document.Content)
	}
	if document.DocumentName != "testdoc" {
		t.Errorf("document name is incorrect: %v", document.DocumentName)
	}
	if document.UserID != "testuser" {
		t.Errorf("document owner is incorrect: %v", document.UserID)
	}
	//
	// // Test deleting a document
	// err = client.DeleteDocument("testuser", "testdoc")
	// if err != nil {
	// 	t.Errorf("failed to delete document: %v", err)
	// }
}
