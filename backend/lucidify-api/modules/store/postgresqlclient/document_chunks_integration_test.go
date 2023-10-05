// //go:build integration
// // +build integration
package postgresqlclient

import (
	"lucidify-api/modules/store/storemodels"
	"testing"
)

func TestChunkFunctions(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

	user := User{
		UserID:           "document_chunks_integration_test_user_id",
		ExternalID:       "TestDocumentChunksID",
		Username:         "TestDocumentChunksUsername",
		PasswordEnabled:  true,
		Email:            "TestDocumentChunks@example.com",
		FirstName:        "TestDocumentChunksFirstName",
		LastName:         "TestDocumentChunksLastName",
		ImageURL:         "https://TestDocumentChunks.com/image.jpg",
		ProfileImageURL:  "https://TestDocumentChunks.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err = store.CreateUserInUsersTable(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	// Create a document first
	doc := &storemodels.Document{
		UserID:       "document_chunks_integration_test_user_id",
		DocumentName: "test_document_name",
		Content:      "test_content",
	}

	insertedDoc, err := store.UploadDocument(doc.UserID, doc.DocumentName, doc.Content)
	if err != nil {
		t.Errorf("Failed to upload test document: %v", err)
	}

	// Create sample chunks
	chunks := []storemodels.Chunk{
		{
			DocumentID:   insertedDoc.DocumentUUID,
			ChunkContent: "chunk_content_1",
			ChunkIndex:   1,
		},
		{
			DocumentID:   insertedDoc.DocumentUUID,
			ChunkContent: "chunk_content_2",
			ChunkIndex:   2,
		},
	}

	// Test UploadChunks
	err = store.UploadChunks(chunks)
	if err != nil {
		t.Errorf("Failed to upload chunks: %v", err)
	}

	// ... rest of your test ...

	t.Cleanup(func() {
		// Cleanup: Delete the chunks and any other resources you've created for the test
		err = store.DeleteAllChunksByDocumentID(insertedDoc.DocumentUUID)
		if err != nil {
			t.Errorf("Failed to delete test chunks: %v", err)
		}

		// Also delete the test document
		err = store.DeleteDocument(doc.UserID, doc.DocumentName)
		if err != nil {
			t.Errorf("Failed to delete test document: %v", err)
		}

		err = store.DeleteUserInUsersTable("document_chunks_integration_test_user_id")
		if err != nil {
			t.Errorf("Failed to delete test user: %v", err)
		}

	})
}
