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
	UploadDocument(documentID, userID, name, content string) error
	GetDocument(documentID string) (*Document, error)
	UpdateDocumentContent(documentID, content string) error
	UpdateDocumentName(documentID, documentName string) error
	DeleteDocument(documentID string) error
}

type WeaviateClientImpl struct {
	client *weaviate.Client
}

type Document struct {
	UserID       string `json:"userId"`
	DocumentName string `json:"documentName"`
	Content      string `json:"content"`
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

	if !classExists(client, "Documents") {
		createWeaviateDocumentsClass(client)
	}

	return &WeaviateClientImpl{client: client}, nil
}

func (w *WeaviateClientImpl) GetWeaviateClient() *weaviate.Client {
	return w.client
}

func (w *WeaviateClientImpl) UploadDocument(documentID, userID, name, content string) error {
	document := map[string]interface{}{
		"documentId":   documentID,
		"userId":       userID,
		"documentName": name,
		"content":      content,
	}

	_, err := w.client.Data().Creator().
		WithID(documentID).
		WithClassName("Documents").
		WithProperties(document).
		Do(context.Background())

	return err
}

func (w *WeaviateClientImpl) GetDocument(documentID string) (*Document, error) {
	objects, err := w.client.Data().ObjectsGetter().
		WithClassName("Documents").
		WithID(documentID).
		Do(context.Background())
	if err != nil {
		// handle error
		return nil, err // it's better to return the error rather than panic
	}

	// If no objects are returned, return an error
	if len(objects) == 0 {
		return nil, errors.New("no documents found")
	}

	// Assume the first object is the one you're looking for
	obj := objects[0]

	// Convert the object to a Document
	doc := &Document{
		UserID:       obj.Properties.(map[string]interface{})["userId"].(string),
		DocumentName: obj.Properties.(map[string]interface{})["documentName"].(string),
		Content:      obj.Properties.(map[string]interface{})["content"].(string),
	}

	return doc, nil
}

func (w *WeaviateClientImpl) UpdateDocumentContent(documentID, content string) error {
	document := map[string]interface{}{
		"content": content,
	}

	err := w.client.Data().Updater().
		WithMerge().
		WithID(documentID).
		WithClassName("Documents").
		WithProperties(document).
		Do(context.Background())

	return err
}

func (w *WeaviateClientImpl) UpdateDocumentName(documentID, documentName string) error {
	document := map[string]interface{}{
		"documentName": documentName,
	}

	err := w.client.Data().Updater().
		WithMerge().
		WithID(documentID).
		WithClassName("Documents").
		WithProperties(document).
		Do(context.Background())

	return err
}

func (w *WeaviateClientImpl) DeleteDocument(documentID string) error {
	err := w.client.Data().Deleter().
		WithClassName("Documents").
		WithID(documentID).
		Do(context.Background())

	return err
}

func classExists(client *weaviate.Client, className string) bool {
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

	// Check if the class already exists
	if classExists(client, "Documents") {
		log.Println("Class 'Documents' already exists")
		return
	}

	classObj := &models.Class{
		Class:       "Documents",
		Description: "A document with associated metadata",
		Vectorizer:  "text2vec-openai",
		Properties: []*models.Property{
			{
				DataType:    []string{"string"},
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

//
//func getDocumentID(userID, name string) string {
//	// Implement the logic to get the document ID based on userID and name
//	return "some-document-id"
//}
