// //go:build integration
// // +build integration
package store

import (
	"log"
	"lucidify-api/modules/store/postgresqlclient"
	"os"
	"testing"
)

func createTestUserInDb() string {
	db, err := postgresqlclient.NewPostgreSQL()

	// the user id registered by the jwt token must exist in the local database
	user := postgresqlclient.User{
		UserID:           "TestStoreIntegrationTestUUID",
		ExternalID:       "TestStoreIntegrationTest",
		Username:         "TestStoreIntegrationTest",
		PasswordEnabled:  true,
		Email:            "TestStoreIntegrationTest@example.com",
		FirstName:        "TestStoreIntegrationTestFirstName",
		LastName:         "TestStoreIntegrationTestLastName",
		ImageURL:         "https://TestStoreIntegrationTestURL.com/image.jpg",
		ProfileImageURL:  "https://TestStoreIntegrationTestProfileURL.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	db.DeleteUserInUsersTable(user.UserID)
	err = db.CheckUserDeletedInUsersTable(user.UserID, 3)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	err = db.CreateUserInUsersTable(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Check if the user exists
	err = db.CheckIfUserInUsersTable(user.UserID, 3)
	if err != nil {
		log.Fatalf("User not found after creation: %v", err)
	}
	return user.UserID
}

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
		{"test_doc_user1_01.txt", 4},
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

func TestUploadDocumentIntegration(t *testing.T) {
	// postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	// if err != nil {
	// 	t.Fatalf("failed to initialize PostgreSQL client: %v", err)
	// }
	//
	// // Initialize the Weaviate client
	// weaviateDB, err := weaviateclient.NewWeaviateClient()
	// if err != nil {
	// 	t.Fatalf("failed to initialize Weaviate client: %v", err)
	// }
	//
	// // Initialize the DocumentService
	// documentService := NewDocumentService(postgresqlDB, weaviateDB)
	//
	// // Define test document parameters
	// userID := createTestUserInDb()
	// name := "test-document-name"
	// content := "test-document-content"
	//
	// // Attempt to upload the document
	// document, err := documentService.UploadDocument(userID, name, content)
	// if err != nil {
	// 	t.Fatalf("failed to upload document: %v", err)
	// }
	// log.Printf("document: %+v", document)

	// Verify that the document was uploaded to PostgreSQL
	// doc, err := postgresqlDB.GetDocument(userID, name)
	// if err != nil || doc == nil {
	// 	t.Fatalf("failed to retrieve document from PostgreSQL: %v", err)
	// }

}

// 	// Verify that the document was uploaded to Weaviate
// 	doc2, err := weaviateDB.GetDocument(document.DocumentUUID.String())
// 	if err != nil || doc2 == nil {
// 		t.Fatalf("failed to retrieve document from Weaviate: %v", err)
// 	}
//
// 	// Clean up: delete the uploaded document and user
// 	t.Cleanup(func() {
// 		err := documentService.DeleteDocument(userID, name, document.DocumentUUID.String())
// 		if err != nil {
// 			t.Errorf("failed to delete test document: %v", err)
// 		}
// 		err = postgresqlDB.DeleteUserInUsersTable(userID)
// 		if err != nil {
// 			t.Errorf("failed to delete test user: %v", err)
// 		}
// 	})
// }
//
// func TestUpdateDocumentNameIntegration(t *testing.T) {
// 	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
// 	if err != nil {
// 		t.Fatalf("failed to initialize PostgreSQL client: %v", err)
// 	}
//
// 	// Initialize the Weaviate client
// 	weaviateDB, err := weaviateclient.NewWeaviateClient()
// 	if err != nil {
// 		t.Fatalf("failed to initialize Weaviate client: %v", err)
// 	}
//
// 	// Initialize the DocumentService
// 	documentService := NewDocumentService(postgresqlDB, weaviateDB)
//
// 	// Define test document parameters
// 	userID := createTestUserInDb()
// 	name := "test-document-name"
// 	content := "test-document-content"
// 	newName := "updated-document-name"
//
// 	// Attempt to upload the document
// 	document, err := documentService.UploadDocument(userID, name, content)
// 	if err != nil {
// 		t.Fatalf("failed to upload document: %v", err)
// 	}
//
// 	// Attempt to update the document name
// 	err = documentService.UpdateDocumentName(document.DocumentUUID.String(), newName)
// 	if err != nil {
// 		t.Fatalf("failed to update document name: %v", err)
// 	}
//
// 	// Verify that the document name was updated in PostgreSQL
// 	doc, err := postgresqlDB.GetDocument(userID, newName)
// 	if err != nil || doc == nil || doc.DocumentName != newName {
// 		t.Fatalf("failed to retrieve document with updated name from PostgreSQL: %v", err)
// 	}
//
// 	// Verify that the document name was updated in Weaviate
// 	doc2, err := weaviateDB.GetDocument(document.DocumentUUID.String())
// 	if err != nil || doc2 == nil || doc2.DocumentName != newName {
// 		t.Fatalf("failed to retrieve document with updated name from Weaviate: %v", err)
// 	}
//
// 	// Clean up: delete the uploaded document and user
// 	t.Cleanup(func() {
// 		err := documentService.DeleteDocument(userID, newName, document.DocumentUUID.String())
// 		if err != nil {
// 			t.Errorf("failed to delete test document: %v", err)
// 		}
// 		err = postgresqlDB.DeleteUserInUsersTable(userID)
// 		if err != nil {
// 			t.Errorf("failed to delete test user: %v", err)
// 		}
// 	})
// }
//
// func TestUpdateDocumentContentIntegration(t *testing.T) {
// 	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
// 	if err != nil {
// 		t.Fatalf("failed to initialize PostgreSQL client: %v", err)
// 	}
//
// 	// Initialize the Weaviate client
// 	weaviateDB, err := weaviateclient.NewWeaviateClient()
// 	if err != nil {
// 		t.Fatalf("failed to initialize Weaviate client: %v", err)
// 	}
//
// 	// Initialize the DocumentService
// 	documentService := NewDocumentService(postgresqlDB, weaviateDB)
//
// 	// Define test document parameters
// 	userID := createTestUserInDb()
// 	name := "test-document-name"
// 	content := "test-document-content"
// 	newContent := "updated-document-content"
//
// 	// Attempt to upload the document
// 	document, err := documentService.UploadDocument(userID, name, content)
// 	if err != nil {
// 		t.Fatalf("failed to upload document: %v", err)
// 	}
//
// 	// Attempt to update the document content
// 	err = documentService.UpdateDocumentContent(document.DocumentUUID.String(), newContent)
// 	if err != nil {
// 		t.Fatalf("failed to update document content: %v", err)
// 	}
//
// 	// Verify that the document content was updated in PostgreSQL
// 	doc, err := postgresqlDB.GetDocument(userID, name)
// 	if err != nil || doc == nil || doc.Content != newContent {
// 		t.Fatalf("failed to retrieve document with updated content from PostgreSQL: %v", err)
// 	}
//
// 	// Verify that the document content was updated in Weaviate
// 	doc2, err := weaviateDB.GetDocument(document.DocumentUUID.String())
// 	if err != nil || doc2 == nil || doc2.Content != newContent {
// 		t.Fatalf("failed to retrieve document with updated content from Weaviate: %v", err)
// 	}
//
// 	// Clean up: delete the uploaded document and user
// 	t.Cleanup(func() {
// 		err := documentService.DeleteDocument(userID, name, document.DocumentUUID.String())
// 		if err != nil {
// 			t.Errorf("failed to delete test document: %v", err)
// 		}
// 		err = postgresqlDB.DeleteUserInUsersTable(userID)
// 		if err != nil {
// 			t.Errorf("failed to delete test user: %v", err)
// 		}
// 	})
// }
