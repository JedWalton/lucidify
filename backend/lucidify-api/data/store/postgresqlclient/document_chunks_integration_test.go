// //go:build integration
// // +build integration
package postgresqlclient

import (
	storemodels2 "lucidify-api/data/store/storemodels"
	"testing"
)

func TestChunkFunctions(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Fatalf("Failed to create test postgresqlclient: %v", err)
	}

	user := storemodels2.User{
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
		t.Fatalf("Failed to create user: %v", err)
	}

	doc := &storemodels2.Document{
		UserID:       user.UserID,
		DocumentName: "test_document_name",
		Content:      "test_content",
	}

	insertedDoc, err := store.UploadDocument(doc.UserID, doc.DocumentName, doc.Content)
	if err != nil {
		t.Fatalf("Failed to upload test document: %v", err)
	}

	chunks := []storemodels2.Chunk{
		{
			UserID:       user.UserID,
			DocumentID:   insertedDoc.DocumentUUID,
			ChunkContent: "chunk_content_1",
			ChunkIndex:   1,
		},
		{
			UserID:       user.UserID,
			DocumentID:   insertedDoc.DocumentUUID,
			ChunkContent: "chunk_content_2",
			ChunkIndex:   2,
		},
	}

	uploadedChunks, err := store.UploadChunks(chunks)
	if err != nil {
		t.Fatalf("Failed to upload chunks: %v", err)
	}

	uploadedChunk1 := uploadedChunks[0]
	if uploadedChunk1.ChunkID.String() == "00000000-0000-0000-0000-000000000000" {
		t.Errorf("Expected chunk ID to be set, but got %s", uploadedChunk1.ChunkID.String())
	}

	err = store.DeleteAllChunksByDocumentID(insertedDoc.DocumentUUID)
	if err != nil {
		t.Fatalf("Failed to delete chunks by document ID: %v", err)
	}

	retrievedChunks, err := store.GetChunksOfDocument(insertedDoc)
	if err != nil {
		t.Fatalf("Failed to retrieve chunks by document ID: %v", err)
	}
	if len(retrievedChunks) != 0 {
		t.Errorf("Expected no chunks, but got %d", len(retrievedChunks))
	}
	retrievedChunksByDocumentID, err := store.GetChunksOfDocumentByDocumentID(insertedDoc.DocumentUUID)
	if err != nil {
		t.Fatalf("Failed to retrieve chunks by document ID: %v", err)
	}
	for i, chunk := range retrievedChunksByDocumentID {
		if chunk != retrievedChunks[i] {
			t.Errorf("Expected chunks to be equal, but got %v and %v", chunk, retrievedChunks[i])
		}
	}

	t.Cleanup(func() {
		// Note, deleting user in users table will delete all associated records
		// of documents and chunks.
		err = store.DeleteAllChunksByDocumentID(insertedDoc.DocumentUUID)
		if err != nil {
			t.Errorf("Failed to delete test chunks: %v", err)
		}

		err = store.DeleteDocument(doc.UserID, doc.DocumentName)
		if err != nil {
			t.Errorf("Failed to delete test document: %v", err)
		}

		err = store.DeleteUserInUsersTable(user.UserID)
		if err != nil {
			t.Errorf("Failed to delete test user: %v", err)
		}
	})
}
