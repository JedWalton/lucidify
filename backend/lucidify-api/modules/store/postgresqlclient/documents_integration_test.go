// //go:build integration
// // +build integration
package postgresqlclient

import (
	"testing"
)

func TestStoreFunctions(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

	user := User{
		UserID:           "documents_integration_test_user_id",
		ExternalID:       "TestCreateUserInUsersTableExternalID",
		Username:         "TestDocumentsIntegrationCreateUserInUsersTableUsername",
		PasswordEnabled:  true,
		Email:            "TestDocumentsCreateUserInUsersTable@example.com",
		FirstName:        "TestDocumentsCreateUserInUsersTableCreateTest",
		LastName:         "TestDocumentsCreateUserInUsersTableUser",
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
	t.Logf("Uploaded document: %+v", doc)

	// Test GetDocument
	docGet, err := store.GetDocument("documents_integration_test_user_id", "test_doc")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}
	if docGet.Content != "test_content" {
		t.Errorf("Expected content 'test_content', got '%s'", docGet.Content)
	}

	// Test GetDocumentByUUID
	documentUUID := doc.DocumentUUID
	docByUUID, err := store.GetDocumentByUUID(documentUUID)
	if err != nil {
		t.Errorf("Failed to get document by UUID: %v", err)
	}
	if docByUUID.DocumentUUID != documentUUID {
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

	// Test DeleteDocumentByUUID
	err = store.DeleteDocumentByUUID(updatedDoc.DocumentUUID)
	if err != nil {
		t.Errorf("Failed to delete document by UUID: %v", err)
	}

	// Verify that the document was deleted
	docByUUID, err = store.GetDocumentByUUID(documentUUID)
	if err == nil || docByUUID != nil {
		t.Errorf("Document should have been deleted, but was still retrievable by UUID")
	}

	t.Cleanup(func() {
		// Delete the test document
		err = store.DeleteDocument("documents_integration_test_user_id", "test_doc")
		if err != nil {
			t.Errorf("Failed to delete test document: %v", err)
		}

		// Delete the test user
		err = store.DeleteUserInUsersTable("documents_integration_test_user_id")
		if err != nil {
			t.Errorf("Failed to delete test user: %v", err)
		}
	})
}
