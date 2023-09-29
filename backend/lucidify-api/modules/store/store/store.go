package store

import (
	"fmt"
	"log"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"
	"time"
)

type DocumentService interface {
	UploadDocument(userID, name, content string) (*postgresqlclient.Document, error)
	GetDocument(userID, name string) (*postgresqlclient.Document, error)
	// GetAllDocuments(userID string) ([]postgresqlclient.Document, error)
	// UpdateDocumentContent(userID, name, content string) error
	DeleteDocument(userID, name, documentUUID string) error
}

type DocumentServiceImpl struct {
	postgresqlDB postgresqlclient.PostgreSQL
	weaviateDB   weaviateclient.WeaviateClient
}

func NewDocumentService(postgresqlDB *postgresqlclient.PostgreSQL, weaviateDB weaviateclient.WeaviateClient) DocumentService {
	return &DocumentServiceImpl{postgresqlDB: *postgresqlDB, weaviateDB: weaviateDB}
}

// func (d *DocumentServiceImpl) UploadDocument(userID, name, content string) error {
// 	document_uuid, err := d.postgresqlDB.UploadDocument(userID, name, content)
// 	if err != nil {
// 		return err
// 	}
// 	err = d.weaviateDB.UploadDocument(document_uuid.String(), userID, name, content)
// 	if err != nil {
// 		err = d.postgresqlDB.DeleteDocument(userID, name)
// 		return err
// 	}
// 	return nil
// }

func (d *DocumentServiceImpl) UploadDocument(userID, name, content string) (*postgresqlclient.Document, error) {
	document, err := d.postgresqlDB.UploadDocument(userID, name, content)
	if err != nil {
		return document, fmt.Errorf("failed to upload document to PostgreSQL: %w", err)
	}

	const maxRetries = 3
	const retryDelay = time.Second * 2
	err = d.weaviateDB.UploadDocument(document.DocumentUUID.String(), userID, name, content)
	if err != nil {
		// Attempt to rollback the PostgreSQL upload.
		var delErr error
		for i := 0; i < maxRetries; i++ {
			delErr = d.postgresqlDB.DeleteDocument(userID, name)
			if delErr == nil {
				break
			}
			// Sleep before retrying
			time.Sleep(retryDelay)
		}

		if delErr != nil {
			return document, fmt.Errorf("failed to upload document to Weaviate: %w; failed to delete document from PostgreSQL after %d retries: %v", err, maxRetries, delErr)
		}
		return document, fmt.Errorf("failed to upload document to Weaviate: %w; document deleted from PostgreSQL", err)
	}

	return document, nil
}

func (d *DocumentServiceImpl) GetDocument(userID, name string) (*postgresqlclient.Document, error) {
	return d.postgresqlDB.GetDocument(userID, name)
}

func (d *DocumentServiceImpl) DeleteDocument(userID, name, documentUUID string) error {
	err := d.postgresqlDB.DeleteDocument(userID, name)
	if err != nil {
		log.Printf("Failed to delete document from PostgreSQL: %v", err)
	}
	err = d.weaviateDB.DeleteDocument(documentUUID)
	if err != nil {
		log.Printf("Failed to delete document from Weaviate: %v", err)
		return err
	}
	return nil
}
