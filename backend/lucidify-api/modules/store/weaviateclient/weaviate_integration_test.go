// //go:build integration
// // +build integration
package weaviateclient

import (
	"fmt"
	"os"
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
	err = weaviateClient.UpdateDocument(documentID, "testuser", "testdoc", "updated test content")
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
	err = weaviateClient.UpdateDocument(documentID, "testuser", "updated testdoc name", "updated test content")
	if err != nil {
		t.Errorf("failed to update document: %v", err)
	}
	document, err = weaviateClient.GetDocument(documentID)
	if document.Content != "updated test content" {
		t.Errorf("document content is incorrect: %v", document.Content)
	}
	if document.DocumentName != "updated testdoc name" {
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

// Helper function to read the content of a file and return it as a string
func readFileContent(filename string) (string, error) {
	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(contentBytes), nil
}

func TestSplitContentIntoChunks(t *testing.T) {
	// Define a struct for test cases
	type testCase struct {
		filename       string
		expectedChunks int
	}

	// Create a slice of test cases
	testCases := []testCase{
		{"test_doc_dogs.txt", 4},
		{"test_doc_cats.txt", 4},
		{"test_doc_vector_databases.txt", 4},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			content, err := readFileContent(tc.filename)
			if err != nil {
				t.Errorf("failed to read file content: %v", err)
			}

			// Use the function to split the content
			chunks, err := splitContentIntoChunks(content)
			if err != nil {
				t.Errorf("failed to split content: %v", err)
			}
			if len(chunks) != tc.expectedChunks {
				t.Errorf("incorrect number of chunks: got %v, want %v", len(chunks), tc.expectedChunks)
			}
		})
	}
}

func setupDocuments(client WeaviateClient, users []string, categories map[string][]string) ([]string, error) {
	var documentIDs []string

	for _, user := range users {
		for i, category := range categories[user] {
			documentID := uuid.New().String()
			documentIDs = append(documentIDs, documentID)
			if err := client.UploadDocument(documentID, user, fmt.Sprintf("testdoc%d", i+1), category); err != nil {
				return nil, fmt.Errorf("failed to upload document: %v", err)
			}
		}
	}

	return documentIDs, nil
}

func teardownDocuments(client WeaviateClient, documentIDs []string) error {
	for _, id := range documentIDs {
		if err := client.DeleteDocument(id); err != nil {
			return fmt.Errorf("failed to delete document with ID %s: %v", id, err)
		}
	}
	return nil
}

func TestSearchDocumentsByText(t *testing.T) {
	weaviateClient, err := NewWeaviateClient()
	if err != nil {
		t.Fatalf("failed to create weaviate client: %v", err)
	}

	users := []string{"testuser1", "testuser2", "testuser3"}
	documents := map[string][]string{
		"testuser1": {
			`Put your first custom data for Cats here.`,
			`Put your second custom data for Cats here.`,
			`Put your third custom data for Cats here.`,
			`Put your fourth custom data for Cats here.`,
			`Put your fifth custom data for Cats here.`,
		},
		"testuser2": {
			`Put your first custom data for Dogs here.`,
			`Put your second custom data for Dogs here.`,
			`Put your third custom data for Dogs here.`,
			`Put your fourth custom data for Dogs here.`,
			`Put your fifth custom data for Dogs here.`,
		},
		"testuser3": {
			`Put your first custom data for Vector Databases here.`,
			`Put your second custom data for Vector Databases here.`,
			`Put your third custom data for Vector Databases here.`,
			`Put your fourth custom data for Vector Databases here.`,
			`Put your fifth custom data for Vector Databases here.`,
		},
	}

	// Keep track of uploaded document IDs for cleanup
	documentIDs, err := setupDocuments(weaviateClient, users, documents)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	defer func() {
		if err := teardownDocuments(weaviateClient, documentIDs); err != nil {
			t.Errorf("teardown failed: %v", err)
		}
	}()

}

//
// 	// Define a query and limit for the test
// 	top_k := 3
// 	userID := "testuser1"
//
// 	concepts := []string{"small animal that goes meow sometimes"}
// 	// Call the SearchDocumentsByText function
// 	result, err := weaviateClient.SearchDocumentsByText(top_k, userID, concepts)
//
// 	if result != nil && result.Data != nil {
// 		getData, ok := result.Data["Get"].(map[string]interface{})
// 		if !ok {
// 			t.Fatalf("unexpected format for 'Get' data")
// 		}
//
// 		documents, ok := getData["Documents"].([]interface{})
// 		if !ok {
// 			t.Fatalf("unexpected format for 'Documents' data")
// 		}
//
// 		for _, document := range documents {
// 			docMap, ok := document.(map[string]interface{})
// 			if !ok {
// 				t.Fatalf("unexpected format for 'document' data")
// 			}
//
// 			documentName := docMap["documentName"].(string)
// 			content := docMap["content"].(string)
// 			additional := docMap["_additional"].(map[string]interface{})
// 			certainty := additional["certainty"].(float64)
// 			distance := additional["distance"].(float64)
//
// 			fmt.Printf("Document Name: %s\n", documentName)
// 			fmt.Printf("Content: %s\n", content)
// 			fmt.Printf("Certainty: %f\n", certainty)
// 			fmt.Printf("Distance: %f\n", distance)
// 		}
// 	}
//
// 	t.Fatalf("SearchDocumentsByText %v", result)
// }
