package store

import (
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"
)

type DocumentService interface {
	// UploadDocument(userID, name, content string) error
	// GetDocument(userID, name string) (postgresqlclient.Document, error)
	// GetAllDocuments(userID string) ([]postgresqlclient.Document, error)
	// UpdateDocumentContent(userID, name, content string) error
	// DeleteDocument(userID, name string) error
}

type DocumentServiceImpl struct {
	postgresqlDB postgresqlclient.PostgreSQL
	weaviateDB   weaviateclient.WeaviateClient
}

func NewDocumentService(postgresqlDB postgresqlclient.PostgreSQL, weaviateDB weaviateclient.WeaviateClient) DocumentService {
	return &DocumentServiceImpl{postgresqlDB: postgresqlDB, weaviateDB: weaviateDB}
}

func (d *DocumentServiceImpl) UploadDocument(userID, name, content string) error {
	document_uuid, err := d.postgresqlDB.UploadDocument(userID, name, content)
	if err != nil {
		return err
	}
	err = d.weaviateDB.UploadDocument(document_uuid.String(), userID, name, content)
	if err != nil {
		err = d.postgresqlDB.DeleteDocument(userID, name)
		return err
	}
	return nil
}
