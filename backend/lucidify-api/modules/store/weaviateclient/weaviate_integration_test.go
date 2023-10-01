// //go:build integration
// // +build integration
package weaviateclient

import (
	"strings"
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

	// Test updating a document content
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

	// Test updating a document name
	err = weaviateClient.UpdateDocumentName(documentID, "updated testdoc")
	if err != nil {
		t.Errorf("failed to update document: %v", err)
	}
	document, err = weaviateClient.GetDocument(documentID)
	if document.Content != "updated test content" {
		t.Errorf("document content is incorrect: %v", document.Content)
	}
	if document.DocumentName != "updated testdoc" {
		t.Errorf("document name is incorrect: %v", document.DocumentName)
	}
	if document.UserID != "testuser" {
		t.Errorf("document owner is incorrect: %v", document.UserID)
	}

	// Test deleting a document
	err = weaviateClient.DeleteDocument(documentID)
	if err != nil {
		t.Errorf("failed to delete document: %v", err)
	}

	_, err = weaviateClient.GetDocument(documentID)
	if err == nil {
		t.Errorf("document was not deleted")
	}
}

func repeatString(str string, count int) string {
	var repeated strings.Builder
	for i := 0; i < count; i++ {
		repeated.WriteString(str)
	}
	return repeated.String()
}

func TestSearchDocumentsByText(t *testing.T) {
	weaviateClient, err := NewWeaviateClient()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	users := []string{"testuser1", "testuser2", "testuser3"}
	categories := map[string]string{
		"testuser1": "Cats " + repeatString("cat ", 20),                         // Repeating "cat " 20 times to make up around 100 words
		"testuser2": "Dogs " + repeatString("dog ", 20),                         // Repeating "dog " 20 times to make up around 100 words
		"testuser3": "Vector Databases " + repeatString("vector database ", 10), // Repeating "vector database " 10 times to make up around 100 words
	}

	// Upload 20 documents
	for i := 0; i < 20; i++ {
		documentID := uuid.New().String()
		user := users[i%3]           // Rotate between the three users
		category := categories[user] // Get category associated with the user
		err = weaviateClient.UploadDocument(documentID, user, "testdoc", category)
		if err != nil {
			t.Errorf("failed to upload document: %v", err)
		}
	}

	t.Log("Initializing the Weaviate client")
	client, err := NewWeaviateClient()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	// Define a query and limit for the test
	// query := "test"
	limit := 3
	userID := "testuser1"

	// Call the SearchDocumentsByText function
	res, err := client.SearchDocumentsByText(limit, userID)
	// if err != nil {
	// 	t.Fatalf("failed to search documents by text: %v", err)
	// }
	t.Logf("SearchDocumentsByText %v", res)

	t.Fatalf("SearchDocumentsByText %v", res)
	// Validate the results
	// For example, check if the returned slice is not nil

	// Additional validations can be added here
	// For example, checking the length of the returned slice, checking the content of the returned documents, etc.
}
