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
