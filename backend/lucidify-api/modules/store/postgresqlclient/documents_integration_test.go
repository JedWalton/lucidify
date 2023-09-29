// //go:build integration
// // +build integration
package postgresqlclient

import (
	"lucidify-api/modules/config"
	"testing"
)

func TestStoreFunctions(t *testing.T) {
	testconfig := config.NewServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewPostgreSQL(PostgresqlURL)
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

	user := User{
		UserID:           "documents_integration_test_user_id",
		ExternalID:       "TestCreateUserInUsersTableExternalID",
		Username:         "TestDocumentsCreateUserInUsersTableUsername",
		PasswordEnabled:  true,
		Email:            "TestCreateUserInUsersTable@example.com",
		FirstName:        "TestCreateUserInUsersTableCreateTest",
		LastName:         "TestCreateUserInUsersTableUser",
		ImageURL:         "https://TestCreateUserInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://TestCreateUserInUsersTable.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err = store.CreateUserInUsersTable(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	// Test UploadDocument
	doc, err := store.UploadDocument("documents_integration_test_user_id", "test_doc", "test_content")
	if err != nil {
		t.Errorf("Failed to upload document: %v", err)
	}

	// Test GetDocument
	docGet, err := store.GetDocument("documents_integration_test_user_id", "test_doc")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}
	if docGet.Content != "test_content" {
		t.Errorf("Expected content 'test_content', got '%s'", docGet.Content)
	}

	// Test GetDocumentByUUID
	documentUUID := doc.DocumentUUID.String()
	docByUUID, err := store.GetDocumentByUUID(documentUUID)
	if err != nil {
		t.Errorf("Failed to get document by UUID: %v", err)
	}
	if docByUUID.DocumentUUID.String() != documentUUID {
		t.Errorf("Expected UUID '%s', got '%s'", documentUUID, docByUUID.DocumentUUID)
	}
	if docByUUID.Content != "test_content" {
		t.Errorf("Expected content 'test_content', got '%s'", docByUUID.Content)
	}

	// Test UpdateDocument
	err = store.UpdateDocument("documents_integration_test_user_id", "test_doc", "updated_content")
	if err != nil {
		t.Errorf("Failed to update document: %v", err)
	}

	updatedDoc, err := store.GetDocument("documents_integration_test_user_id", "test_doc")
	if err != nil {
		t.Errorf("Failed to get updated document: %v", err)
	}
	if updatedDoc.Content != "updated_content" {
		t.Errorf("Expected updated content 'updated_content', got '%s'", updatedDoc.Content)
	}

	// Test GetAllDocuments
	docs, err := store.GetAllDocuments("documents_integration_test_user_id")
	if err != nil {
		t.Errorf("Failed to get all documents: %v", err)
	}
	if len(docs) != 1 {
		t.Errorf("Expected 1 document, got %d", len(docs))
	}

	// Test UpdateDocumentName
	newDocumentName := "new_test_doc"
	err = store.UpdateDocumentName(doc.DocumentUUID, newDocumentName)
	if err != nil {
		t.Errorf("Failed to update document name: %v", err)
	}

	// Verify that the document name was updated
	updatedDoc, err = store.GetDocument("documents_integration_test_user_id", newDocumentName)
	if err != nil {
		t.Errorf("Failed to get document with new name: %v", err)
	}
	if updatedDoc.DocumentName != newDocumentName {
		t.Errorf("Expected updated document name '%s', got '%s'", newDocumentName, updatedDoc.DocumentName)
	}

	// Test UpdateDocumentContent
	newContent := "new_updated_content"
	err = store.UpdateDocumentContent(updatedDoc.DocumentUUID, newContent)
	if err != nil {
		t.Errorf("Failed to update document content: %v", err)
	}

	// Verify that the document content was updated
	updatedDoc, err = store.GetDocument("documents_integration_test_user_id", newDocumentName)
	if err != nil {
		t.Errorf("Failed to get document with new content: %v", err)
	}
	if updatedDoc.Content != newContent {
		t.Errorf("Expected updated content '%s', got '%s'", newContent, updatedDoc.Content)
	}

	t.Cleanup(func() {
		err = store.DeleteUserInUsersTable("documents_integration_test_user_id")
		if err != nil {
			t.Errorf("Failed to delete test user: %v", err)
		}

		err = store.DeleteDocument("documents_integration_test_user_id", "test_doc")
		if err != nil {
			t.Errorf("Failed to delete test document: %v", err)
		}
	})
}
