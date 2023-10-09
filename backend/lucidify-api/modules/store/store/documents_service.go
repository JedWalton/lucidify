package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/storemodels"
	"lucidify-api/modules/store/weaviateclient"
	"net/http"

	"github.com/google/uuid"
)

type DocumentService interface {
	UploadDocument(userID, name, content string) (*storemodels.Document, error)
	GetDocument(userID, name string) (*storemodels.Document, error)
	GetAllDocuments(userID string) ([]storemodels.Document, error)
	DeleteDocument(documentID uuid.UUID) error
	UpdateDocumentName(documentID uuid.UUID, name string) error
	UpdateDocumentContent(documentUUID uuid.UUID, content string) error
}

type DocumentServiceImpl struct {
	postgresqlDB postgresqlclient.PostgreSQL
	weaviateDB   weaviateclient.WeaviateClient
}

func NewDocumentService(
	postgresqlDB *postgresqlclient.PostgreSQL,
	weaviateDB weaviateclient.WeaviateClient) DocumentService {
	return &DocumentServiceImpl{postgresqlDB: *postgresqlDB, weaviateDB: weaviateDB}
}

func (d *DocumentServiceImpl) UploadDocument(
	userID, name, content string) (*storemodels.Document, error) {

	var cleanupTasks []func() error
	shouldCleanup := false

	defer func() {
		if shouldCleanup {
			for _, task := range cleanupTasks {
				if err := task(); err != nil {
					// Log the cleanup error or handle it as needed
					log.Printf("Failed to cleanup: %v", err)
				}
			}
		}
	}()

	document, err := d.postgresqlDB.UploadDocument(userID, name, content)
	if err != nil {
		// Append a cleanup task for the uploaded document in PostgreSQL
		cleanupTasks = append(cleanupTasks, func() error {
			return d.postgresqlDB.DeleteDocumentByUUID(document.DocumentUUID)
		})
		shouldCleanup = true
		return document, fmt.Errorf("Upload failed at upload document to PostgreSQL: %w", err)
	}

	chunks, err := splitContentIntoChunks(*document)
	if err != nil {
		return document, fmt.Errorf("Upload failed at split content into chunks: %w", err)
	}

	chunksWithChunkID, err := d.postgresqlDB.UploadChunks(chunks)
	if err != nil {
		cleanupTasks = append(cleanupTasks, func() error {
			return d.postgresqlDB.DeleteDocumentByUUID(document.DocumentUUID)
		})
		shouldCleanup = true
		return document, fmt.Errorf("Upload failed at upload chunks to PostgreSQL: %w", err)
	}

	err = d.weaviateDB.UploadChunks(chunksWithChunkID)
	if err != nil {
		shouldCleanup = true
		cleanupTasks = append(cleanupTasks, func() error {
			return d.DeleteDocument(document.DocumentUUID)
		})
		return document, fmt.Errorf("Upload failed at upload chunks to weaviate: %w", err)
	}

	return document, nil
}

func splitContentIntoChunks(document storemodels.Document) ([]storemodels.Chunk, error) {
	cfg := config.NewServerConfig()

	url := cfg.AI_API_URL + "/split_text_to_chunks"
	payload := map[string]string{
		"text": document.Content,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AI-API-KEY", cfg.X_AI_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, body)
	}

	var chunkContents []string
	if err := json.NewDecoder(resp.Body).Decode(&chunkContents); err != nil {
		return nil, err
	}

	var chunks []storemodels.Chunk
	for index, content := range chunkContents {
		chunk := storemodels.Chunk{
			UserID:       document.UserID,
			DocumentID:   document.DocumentUUID,
			ChunkContent: content,
			ChunkIndex:   index,
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func (d *DocumentServiceImpl) GetDocument(userID, name string) (*storemodels.Document, error) {
	return d.postgresqlDB.GetDocument(userID, name)
}

func (d *DocumentServiceImpl) GetAllDocuments(userID string) ([]storemodels.Document, error) {
	return d.postgresqlDB.GetAllDocuments(userID)
}

func (d *DocumentServiceImpl) DeleteDocument(documentUUID uuid.UUID) error {
	chunks, err := d.postgresqlDB.GetChunksOfDocumentByDocumentID(documentUUID)
	if err != nil {
		return fmt.Errorf("Failed to get chunks of document: %w", err)
	}
	err = d.weaviateDB.DeleteChunks(chunks)
	if err != nil {
		return fmt.Errorf("Failed to delete chunks from Weaviate: %w", err)
	}
	err = d.postgresqlDB.DeleteDocumentByUUID(documentUUID)
	if err != nil {
		log.Printf("Failed to delete document from PostgreSQL: %v", err)
	}
	return nil
}

func (d *DocumentServiceImpl) UpdateDocumentName(documentID uuid.UUID, name string) error {
	err := d.postgresqlDB.UpdateDocumentName(documentID, name)
	if err != nil {
		return fmt.Errorf("Failed to update document name in PostgreSQL: %w", err)
	}

	return nil
}

func (d *DocumentServiceImpl) UpdateDocumentContent(documentID uuid.UUID, content string) error {
	// First, get the document by UUID to ensure it exists and to get the current content.
	documentBeforeChange, err := d.postgresqlDB.GetDocumentByUUID(documentID)
	if err != nil {
		return fmt.Errorf("Failed to get document by UUID from PostgreSQL: %w", err)
	}

	// Delete the document chunks from PostgreSQL and Weaviate.
	chunksPostgreSQL, err := d.postgresqlDB.GetChunksOfDocumentByDocumentID(documentID)
	if err != nil {
		return fmt.Errorf("Failed to get chunks of document from PostgreSQL: %w", err)
	}

	err = d.weaviateDB.DeleteChunks(chunksPostgreSQL)
	if err != nil {
		return fmt.Errorf("Failed to delete chunks from Weaviate: %w", err)
	}

	if err := d.postgresqlDB.DeleteAllChunksByDocumentID(documentID); err != nil {
		return fmt.Errorf("Failed to delete chunks from PostgreSQL: %w", err)
	}

	// Split the new content into chunks.
	documentWithNewContent := documentBeforeChange
	documentWithNewContent.Content = content
	splitChunksOfNewContent, err := splitContentIntoChunks(*documentWithNewContent)
	if err != nil {
		return fmt.Errorf("Failed to split new content into chunks: %w", err)
	}

	// Upload the new chunks to PostgreSQL and Weaviate.
	chunksWithChunkID, err := d.postgresqlDB.UploadChunks(splitChunksOfNewContent)
	if err != nil {
		return fmt.Errorf("Failed to upload new chunks to PostgreSQL: %w", err)
	}

	err = d.weaviateDB.UploadChunks(chunksWithChunkID)
	if err != nil {
		return fmt.Errorf("Failed to upload new chunks to Weaviate: %w", err)
	}

	return nil
}
