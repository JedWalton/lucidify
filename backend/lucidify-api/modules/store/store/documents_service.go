package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/storemodels"
	"lucidify-api/modules/store/weaviateclient"
	"net/http"
)

type DocumentService interface {
	UploadDocument(userID, name, content string) (*storemodels.Document, error)
	GetDocument(userID, name string) (*storemodels.Document, error)
	GetAllDocuments(userID string) ([]storemodels.Document, error)
	// DeleteDocument(userID, name, documentUUID string) error
	// UpdateDocumentName(documentUUID, name string) error
	// UpdateDocumentContent(documentUUID, content string) error
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

// func (d *DocumentServiceImpl) UploadDocument(
//
//		userID, name, content string) (*storemodels.Document, []func() error, error) {
//
//		var cleanupTasks []func() error
//
//		document, err := d.postgresqlDB.UploadDocument(userID, name, content)
//		if err != nil {
//			return document, cleanupTasks, fmt.Errorf("Upload failed at upload document to PostgreSQL: %w", err)
//		}
//
//		chunks, err := splitContentIntoChunks(*document)
//		if err != nil {
//			return document, cleanupTasks, fmt.Errorf("Upload failed at split content into chunks: %w", err)
//		}
//
//		chunksWithChunkID, err := d.postgresqlDB.UploadChunks(chunks)
//		if err != nil {
//			return document, cleanupTasks, fmt.Errorf("Upload failed at upload chunks to PostgreSQL: %w", err)
//		}
//
//		// Append a cleanup task for the uploaded chunks in PostgreSQL
//		cleanupTasks = append(cleanupTasks, func() error {
//			return d.postgresqlDB.DeleteAllChunksByDocumentID(document.DocumentUUID)
//		})
//
//		err = d.weaviateDB.UploadChunks(chunksWithChunkID)
//		if err != nil {
//			return document, cleanupTasks, fmt.Errorf("Upload failed at upload chunks to weaviate: %w", err)
//		}
//
//		// Append a cleanup task for the uploaded chunks in Weaviate
//		cleanupTasks = append(cleanupTasks, func() error {
//			return d.weaviateDB.DeleteChunks(chunksWithChunkID)
//		})
//
//		return document, cleanupTasks, nil
//	}
// func (d *DocumentServiceImpl) UploadDocument(
// 	userID, name, content string) (*storemodels.Document, error) {
//
// 	var cleanupTasks []func() error
// 	var cleanupRequired bool
// 	defer func() {
// 		if cleanupRequired {
// 			for _, task := range cleanupTasks {
// 				if err := task(); err != nil {
// 					// Log the cleanup error or handle it as needed
// 				}
// 			}
// 		}
// 	}()
//
// 	document, err := d.postgresqlDB.UploadDocument(userID, name, content)
// 	if err != nil {
// 		cleanupRequired = true
// 		return document, fmt.Errorf("Upload failed at upload document to PostgreSQL: %w", err)
// 	}
//
// 	chunks, err := splitContentIntoChunks(*document)
// 	if err != nil {
// 		cleanupRequired = true
// 		return document, fmt.Errorf("Upload failed at split content into chunks: %w", err)
// 	}
//
// 	chunksWithChunkID, err := d.postgresqlDB.UploadChunks(chunks)
// 	if err != nil {
// 		cleanupRequired = true
// 		return document, fmt.Errorf("Upload failed at upload chunks to PostgreSQL: %w", err)
// 	}
//
// 	// Append a cleanup task for the uploaded chunks in PostgreSQL
// 	cleanupTasks = append(cleanupTasks, func() error {
// 		return d.postgresqlDB.DeleteAllChunksByDocumentID(document.DocumentUUID)
// 	})
//
// 	err = d.weaviateDB.UploadChunks(chunksWithChunkID)
// 	if err != nil {
// 		cleanupRequired = true
// 		return document, fmt.Errorf("Upload failed at upload chunks to weaviate: %w", err)
// 	}
//
// 	// Append a cleanup task for the uploaded chunks in Weaviate
// 	cleanupTasks = append(cleanupTasks, func() error {
// 		return d.weaviateDB.DeleteChunks(chunksWithChunkID)
// 	})
//
// 	return document, nil
// }

func (d *DocumentServiceImpl) UploadDocument(
	userID, name, content string) (*storemodels.Document, error) {

	var cleanupTasks []func() error
	defer func() {
		for _, task := range cleanupTasks {
			if err := task(); err != nil {
				// Log the cleanup error or handle it as needed
			}
		}
	}()

	document, err := d.postgresqlDB.UploadDocument(userID, name, content)
	if err != nil {
		return document, fmt.Errorf("Upload failed at upload document to PostgreSQL: %w", err)
	}
	// Append a cleanup task for the uploaded document in PostgreSQL
	cleanupTasks = append(cleanupTasks, func() error {
		return d.postgresqlDB.DeleteDocumentByUUID(document.DocumentUUID)
	})

	chunks, err := splitContentIntoChunks(*document)
	if err != nil {
		return document, fmt.Errorf("Upload failed at split content into chunks: %w", err)
	}

	chunksWithChunkID, err := d.postgresqlDB.UploadChunks(chunks)
	if err != nil {
		return document, fmt.Errorf("Upload failed at upload chunks to PostgreSQL: %w", err)
	}
	// Append a cleanup task for the uploaded chunks in PostgreSQL
	cleanupTasks = append(cleanupTasks, func() error {
		return d.postgresqlDB.DeleteAllChunksByDocumentID(document.DocumentUUID)
	})

	err = d.weaviateDB.UploadChunks(chunksWithChunkID)
	if err != nil {
		return document, fmt.Errorf("Upload failed at upload chunks to weaviate: %w", err)
	}
	// Append a cleanup task for the uploaded chunks in Weaviate
	cleanupTasks = append(cleanupTasks, func() error {
		return d.weaviateDB.DeleteChunks(chunksWithChunkID)
	})

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

//
// func (d *DocumentServiceImpl) DeleteDocument(userID, name, documentUUID string) error {
// 	err := d.postgresqlDB.DeleteDocument(userID, name)
// 	if err != nil {
// 		log.Printf("Failed to delete document from PostgreSQL: %v", err)
// 	}
// 	err = d.weaviateDB.DeleteDocument(documentUUID)
// 	if err != nil {
// 		return fmt.Errorf("Failed to delete document from Weaviate: %w", err)
// 	}
// 	return nil
// }
//
// func (d *DocumentServiceImpl) UpdateDocumentName(documentUUID, name string) error {
// 	parsedDocumentUUID, err := uuid.Parse(documentUUID)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse UUID: %w", err)
// 	}
//
// 	documentBeforeChange, err := d.postgresqlDB.GetDocumentByUUID(documentUUID)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = d.postgresqlDB.UpdateDocumentName(parsedDocumentUUID, name)
// 	if err != nil {
// 		return fmt.Errorf("failed to update document name in PostgreSQL: %w", err)
// 	}
//
// 	// Try to update the name in Weaviate
// 	err = d.weaviateDB.UpdateDocument(documentUUID, documentBeforeChange.UserID, name, documentBeforeChange.Content)
// 	if err != nil {
// 		// Log the error and try to revert the change in PostgreSQL
// 		log.Printf("Failed to update document name in Weaviate: %v. Returning postgresql name back to original", err)
// 		revertErr := d.postgresqlDB.UpdateDocumentName(parsedDocumentUUID, documentBeforeChange.DocumentName)
// 		if revertErr != nil {
// 			log.Printf("Failed to restore document name to original name in PostgreSQL: %v", revertErr)
// 			// Consider whether to return the original error, the revert error, or both
// 			return fmt.Errorf("failed to update document name in Weaviate, and failed to revert change in PostgreSQL: %w, revert error: %v", err, revertErr)
// 		}
// 		// If revert was successful, return the original error
// 		return fmt.Errorf("failed to update document name in Weaviate: %w", err)
// 	}
//
// 	return nil
// }
//
// func (d *DocumentServiceImpl) UpdateDocumentContent(documentUUID, content string) error {
// 	// First, get the document by UUID to ensure it exists and to get the current content.
// 	documentBeforeChange, err := d.postgresqlDB.GetDocumentByUUID(documentUUID)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Update the content in the PostgreSQL database.
// 	err = d.postgresqlDB.UpdateDocumentContent(uuid.MustParse(documentUUID), content)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Update the content in the Weaviate database.
// 	err = d.weaviateDB.UpdateDocument(documentUUID, documentBeforeChange.UserID, documentBeforeChange.DocumentName, content)
// 	if err != nil {
// 		// If updating in Weaviate fails, rollback the change in PostgreSQL.
// 		log.Printf("Failed to update document content in Weaviate: %v. Returning PostgreSQL content back to original", err)
// 		errRollback := d.postgresqlDB.UpdateDocumentContent(uuid.MustParse(documentUUID), documentBeforeChange.Content)
// 		if errRollback != nil {
// 			log.Printf("Failed to restore document content to original in PostgreSQL: %v", errRollback)
// 		}
// 		return err
// 	}
//
// 	return nil
// }
