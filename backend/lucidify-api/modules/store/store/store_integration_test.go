// //go:build integration
// // +build integration
package store

import (
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"
	"testing"
)

func createTestUserInDb() string {
	testconfig := config.NewServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL
	db, err := postgresqlclient.NewPostgreSQL(PostgresqlURL)

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

func TestUploadDocumentIntegration(t *testing.T) {
	// Initialize the PostgreSQL client
	config := config.NewServerConfig()
	postgresqlURL := config.PostgresqlURL
	postgresqlDB, err := postgresqlclient.NewPostgreSQL(postgresqlURL)
	if err != nil {
		t.Fatalf("failed to initialize PostgreSQL client: %v", err)
	}

	// Initialize the Weaviate client
	weaviateDB, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		t.Fatalf("failed to initialize Weaviate client: %v", err)
	}

	// Initialize the DocumentService
	documentService := NewDocumentService(postgresqlDB, weaviateDB)

	// Define test document parameters
	userID := createTestUserInDb()
	name := "test-document-name"
	content := "test-document-content"

	// Attempt to upload the document
	document, err := documentService.UploadDocument(userID, name, content)
	if err != nil {
		t.Fatalf("failed to upload document: %v", err)
	}

	// Verify that the document was uploaded to PostgreSQL
	doc, err := postgresqlDB.GetDocument(userID, name)
	if err != nil || doc == nil {
		t.Fatalf("failed to retrieve document from PostgreSQL: %v", err)
	}

	// Verify that the document was uploaded to Weaviate
	doc2, err := weaviateDB.GetDocument(document.DocumentUUID.String())
	if err != nil || doc2 == nil {
		t.Fatalf("failed to retrieve document from Weaviate: %v", err)
	}

	// Clean up: delete the uploaded document
	err = documentService.DeleteDocument(userID, name, document.DocumentUUID.String())
	if err != nil {
		t.Fatalf("failed to delete test document: %v", err)
	}
}
