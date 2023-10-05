// //go:build integration
// // +build integration
package weaviateclient

import (
	"testing"
)

func TestWeaviateClient(t *testing.T) {
	// weaviateClient, err := NewWeaviateClient()
	// if err != nil {
	// 	t.Fatalf("failed to create weaviate client: %v", err)
	// }

	// Test uploading a document
	// documentID := uuid.New().String()
	// err = weaviateClient.UploadDocument(documentID, "testuser", "testdoc", "test content")
	// if err != nil {
	// 	t.Errorf("failed to upload document: %v", err)
	// }
}

//		document, err := weaviateClient.GetDocument(documentID)
//		if err != nil {
//			t.Errorf("failed to get document: %v", err)
//		}
//		t.Logf("document: %+v", document)
//		if document.UserID != "testuser" {
//			t.Errorf("document owner is incorrect: %v", document.UserID)
//		}
//		if document.DocumentName != "testdoc" {
//			t.Errorf("document name is incorrect: %v", document.DocumentName)
//		}
//		if document.Content != "test content" {
//			t.Errorf("document content is incorrect: %v", document.Content)
//		}
//
//		// Test updating a document content
//		err = weaviateClient.UpdateDocument(documentID, "testuser", "testdoc", "updated test content")
//		if err != nil {
//			t.Errorf("failed to update document: %v", err)
//		}
//		document, err = weaviateClient.GetDocument(documentID)
//		if document.Content != "updated test content" {
//			t.Errorf("document content is incorrect: %v", document.Content)
//		}
//		if document.DocumentName != "testdoc" {
//			t.Errorf("document name is incorrect: %v", document.DocumentName)
//		}
//		if document.UserID != "testuser" {
//			t.Errorf("document owner is incorrect: %v", document.UserID)
//		}
//
//		// Test updating a document name
//		err = weaviateClient.UpdateDocument(documentID, "testuser", "updated testdoc name", "updated test content")
//		if err != nil {
//			t.Errorf("failed to update document: %v", err)
//		}
//		document, err = weaviateClient.GetDocument(documentID)
//		if document.Content != "updated test content" {
//			t.Errorf("document content is incorrect: %v", document.Content)
//		}
//		if document.DocumentName != "updated testdoc name" {
//			t.Errorf("document name is incorrect: %v", document.DocumentName)
//		}
//		if document.UserID != "testuser" {
//			t.Errorf("document owner is incorrect: %v", document.UserID)
//		}
//
//		// Test deleting a document
//		err = weaviateClient.DeleteDocument(documentID)
//		if err != nil {
//			t.Errorf("failed to delete document: %v", err)
//		}
//
//		_, err = weaviateClient.GetDocument(documentID)
//		if err == nil {
//			t.Errorf("document was not deleted")
//		}
//	}
//
// // Helper function to read the content of a file and return it as a string

//
// func setupDocuments(client WeaviateClient) ([]string, error) {
// 	var documentIDs []string
//
// 	userDocuments := map[string][]string{
// 		"testuser1": {
// 			"test_doc_user1_01.txt",
// 			// "test_doc_testuser1_02.txt",
// 			// "test_doc_testuser1_03.txt",
// 		},
// 		// Add more users and their documents as needed
// 	}
//
// 	for user, docs := range userDocuments {
// 		for _, doc := range docs {
// 			documentID := uuid.New().String()
// 			documentIDs = append(documentIDs, documentID)
//
// 			content, err := readFileContent(doc)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to read file content for %s: %v", doc, err)
// 			}
//
// 			// Assuming the document name in the UploadDocument function is the same as the filename
// 			if err := client.UploadDocument(documentID, user, doc, content); err != nil {
// 				return nil, fmt.Errorf("failed to upload document %s for user %s: %v", doc, user, err)
// 			}
// 		}
// 	}
//
// 	return documentIDs, nil
// }
//
// func teardownDocuments(client WeaviateClient, documentIDs []string) error {
// 	for _, id := range documentIDs {
// 		if err := client.DeleteDocument(id); err != nil {
// 			return fmt.Errorf("failed to delete document with ID %s: %v", id, err)
// 		}
// 	}
// 	return nil
// }
//
// func TestSearchDocumentsByText(t *testing.T) {
// 	weaviateClient, err := NewWeaviateClient()
// 	if err != nil {
// 		t.Fatalf("failed to create weaviate client: %v", err)
// 	}
//
// 	// Keep track of uploaded document IDs for cleanup
// 	documentIDs, err := setupDocuments(weaviateClient)
// 	if err != nil {
// 		t.Fatalf("setup failed: %v", err)
// 	}
//
// 	defer func() {
// 		if err := teardownDocuments(weaviateClient, documentIDs); err != nil {
// 			t.Errorf("teardown failed: %v", err)
// 		}
// 	}()
//
// }

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
