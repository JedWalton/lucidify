package weaviateclient

import (
	"context"
	"errors"
	"log"
	"lucidify-api/modules/config"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

type WeaviateClient interface {
	GetWeaviateClient() *weaviate.Client
	UploadDocument(userID, name, content string) error
	//DeleteDocument(userID, name string) error
	//UpdateDocument(userID, name, content string) error
}

type WeaviateClientImpl struct {
	client *weaviate.Client
}

func NewWeaviateClient() (WeaviateClient, error) {
	config := config.NewServerConfig()
	cfg := weaviate.Config{
		Host:   "localhost:8090",
		Scheme: "http",
		Headers: map[string]string{
			"X-OpenAI-Api-Key": config.OPENAI_API_KEY,
		},
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.New("client is nil after initialization")
	}
	if !doesClassExist(client, "documents") {
		createWeaviateDocumentsClass(client)
	}
	return &WeaviateClientImpl{client: client}, nil
}

func (w *WeaviateClientImpl) GetWeaviateClient() *weaviate.Client {
	return w.client
}

func (w *WeaviateClientImpl) UploadDocument(userID, name, content string) error {
	document := map[string]interface{}{
		"userId":       userID,
		"documentName": name,
		"content":      content,
	}

	_, err := w.client.Data().Creator().
		WithID("postgres_doc_id_here").
		WithClassName("documents").
		WithProperties(document).
		Do(context.Background())

	return err
}

func doesClassExist(client *weaviate.Client, className string) bool {
	if client == nil {
		log.Println("Client is nil in doesClassExist")
		return false
	}
	schema, err := client.Schema().ClassGetter().WithClassName(className).Do(context.Background())
	if err != nil {
		return false
	}
	log.Printf("%v", schema)
	return true
}

func createWeaviateDocumentsClass(client *weaviate.Client) {
	if client == nil {
		log.Println("Client is nil in createWeaviateDocumentsClass")
		return
	}

	classObj := &models.Class{
		Class:       "documents",
		Description: "A document with associated metadata",
		Vectorizer:  "text2vec-openai",
		Properties: []*models.Property{
			{
				DataType:    []string{"int"},
				Description: "Unique identifier of the document",
				Name:        "documentId",
			},
			{
				DataType:    []string{"string"},
				Description: "User identifier associated with the document",
				Name:        "userId",
			},
			{
				DataType:    []string{"string"},
				Description: "Name of the document",
				Name:        "documentName",
			},
			{
				DataType:    []string{"text"},
				Description: "Content of the document",
				Name:        "content",
			},
			{
				DataType:    []string{"date"},
				Description: "Creation timestamp of the document",
				Name:        "createdAt",
			},
			{
				DataType:    []string{"date"},
				Description: "Update timestamp of the document",
				Name:        "updatedAt",
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		panic(err)
	}
}

//func (w *WeaviateClientImpl) DeleteDocument(userID, name string) error {
//	documentID := getDocumentID(userID, name)
//	err := w.client.Data().Deleter().
//		WithClassName("documents").
//		WithID(documentID).
//		Do(context.Background())
//
//	return err
//}
//
//func (w *WeaviateClientImpl) UpdateDocument(userID, name, content string) error {
//	documentID := getDocumentID(userID, name)
//	document := map[string]interface{}{
//		"content": content,
//	}
//
//	err := w.client.Data().Updater().
//		WithClassName("documents").
//		WithID(documentID).
//		WithProperties(document).
//		Do(context.Background())
//
//	return err
//}
//
//func getDocumentID(userID, name string) string {
//	// Implement the logic to get the document ID based on userID and name
//	return "some-document-id"
//}
