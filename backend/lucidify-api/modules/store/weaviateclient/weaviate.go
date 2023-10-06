package weaviateclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/storemodels"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

type WeaviateClient interface {
	GetWeaviateClient() *weaviate.Client
	UploadChunk(storemodels.Chunk) error
	UploadChunks([]storemodels.Chunk) error
	DeleteChunk(chunkID uuid.UUID) error
	DeleteChunks([]storemodels.Chunk) error
	GetChunks(chunksFromPostgresql []storemodels.Chunk) ([]storemodels.Chunk, error)
	// UploadDocument(documentID, userID, name, content string) error
	// GetDocument(documentID string) (*Document, error)
	// UpdateDocument(documentID, userID, name, content string) error
	// DeleteDocument(documentID string) error
	// SearchDocumentsByText(limit int, userID string, concepts []string) (*models.GraphQLResponse, error)
}

type WeaviateClientImpl struct {
	client *weaviate.Client
}

func classExists(client *weaviate.Client, className string) bool {
	schema, err := client.Schema().ClassGetter().WithClassName(className).Do(context.Background())
	if err != nil {
		return false
	}
	log.Printf("%v", schema)
	return true
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
				Description: "Unique identifier of the chunk within the document",
				Name:        "chunkId",
			},
			{
				DataType:    []string{"text"},
				Description: "A chunk of the document content",
				Name:        "chunkContent",
			},
			{
				DataType:    []string{"int"},
				Description: "Index of the chunk in the document",
				Name:        "chunkIndex",
			},
		},
	}

	err := client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
	if err != nil {
		panic(err)
	}
}

func (w *WeaviateClientImpl) UploadChunks(chunks []storemodels.Chunk) error {
	for _, chunk := range chunks {
		err := w.UploadChunk(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *WeaviateClientImpl) UploadChunk(chunk storemodels.Chunk) error {
	if w.client == nil {
		return errors.New("Weaviate client is not initialized")
	}

	// Convert the chunk to a format suitable for Weaviate
	chunkData := map[string]interface{}{
		"documentId":   chunk.DocumentID.String(),
		"userId":       chunk.UserID,
		"chunkId":      chunk.ChunkID.String(), // Convert UUID to string
		"chunkContent": chunk.ChunkContent,
		"chunkIndex":   chunk.ChunkIndex,
	}

	// Use the Weaviate client to upload the chunk
	_, err := w.client.Data().Creator().
		WithID(chunk.ChunkID.String()). // Use chunk's UUID as ID
		WithClassName("Documents").     // Assuming the class name for chunks is "DocumentChunks"
		WithProperties(chunkData).
		Do(context.Background())

	if err != nil {
		return fmt.Errorf("failed to upload chunk: %w", err)
	}

	return nil
}

func (w *WeaviateClientImpl) DeleteAllChunksByDocumentID(documentID string) error {
	return fmt.Errorf("This is not implemented as weaviate does not support non "+
		"index deleting: %w. This func exists in document_service.", errors.New("not implemented"))
}

func (w *WeaviateClientImpl) DeleteChunk(chunkID uuid.UUID) error {
	err := w.client.Data().Deleter().
		WithClassName("Documents").
		WithID(chunkID.String()).
		Do(context.Background())

	return err
}

func (w *WeaviateClientImpl) DeleteChunks(chunks []storemodels.Chunk) error {
	for _, chunk := range chunks {
		err := w.DeleteChunk(chunk.ChunkID)
		if err != nil {
			return err
		}
	}
	return nil
}

//	func (w *WeaviateClientImpl) UploadDocument(documentID, userID, name, content string) error {
//		document := map[string]interface{}{
//			"documentId":   documentID,
//			"userId":       userID,
//			"documentName": name,
//			"content":      content,
//		}
//
//		_, err := w.client.Data().Creator().
//			WithID(documentID).
//			WithClassName("Documents").
//			WithProperties(document).
//			Do(context.Background())
//
//		return err
//	}

// Helper function to split content into chunks
// func splitContentIntoChunks(content string, chunkSize int) []string {
// 	var chunks []string
// 	runes := []rune(content)
//
// 	for i := 0; i < len(runes); i += chunkSize {
// 		end := i + chunkSize
// 		if end > len(runes) {
// 			end = len(runes)
// 		}
// 		chunks = append(chunks, string(runes[i:end]))
// 	}
//
// 	return chunks
// }

// func splitContentIntoChunks(content string) ([]string, error) {
// 	cfg := config.NewServerConfig()
//
// 	url := cfg.AI_API_URL + "/split_text_to_chunks"
// 	payload := map[string]string{
// 		"text": content,
// 	}
// 	jsonPayload, err := json.Marshal(payload)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("X-AI-API-KEY", cfg.X_AI_API_KEY)
//
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
//
// 	if resp.StatusCode != http.StatusOK {
// 		body, _ := io.ReadAll(resp.Body)
// 		return nil, fmt.Errorf("API call failed with status %d: %s", resp.StatusCode, body)
// 	}
//
// 	var chunks []string
// 	if err := json.NewDecoder(resp.Body).Decode(&chunks); err != nil {
// 		return nil, err
// 	}
//
// 	return chunks, nil
// }

//	func (w *WeaviateClientImpl) UploadDocument(documentID, userID, name, content string) error {
//		// func (w *WeaviateClientImpl) UploadDocument(chunks storemodels.Chunk) error {
//
//		chunks, err := splitContentIntoChunks(content)
//		if err != nil {
//			return err
//		}
//
//		for i, chunk := range chunks {
//			document := map[string]interface{}{
//				"documentId": documentID,
//				"userId":     userID,
//				"chunk":      chunk,
//				"chunkId":    i,
//			}
//
//			_, err := w.client.Data().Creator().
//				// WithID(documentID).
//				WithClassName("Documents").
//				WithProperties(document).
//				Do(context.Background())
//
//			if err != nil {
//				return err
//			}
//		}
//
//		return nil
//	}
func (w *WeaviateClientImpl) GetChunks(chunksFromPostgresql []storemodels.Chunk) ([]storemodels.Chunk, error) {
	var chunksFromWeaviate []storemodels.Chunk
	for _, chunk := range chunksFromPostgresql {
		objects, err := w.client.Data().ObjectsGetter().
			WithClassName("Documents").
			WithID(chunk.ChunkID.String()).
			Do(context.Background())

		if err != nil {
			return nil, err
		}

		if len(objects) == 0 {
			return nil, fmt.Errorf("no object found for chunk ID: %s", chunk.ChunkID.String())
		}

		fmt.Printf("objects: %+v\n", objects[0])

		// Extract properties from the first object
		properties := objects[0].Properties.(map[string]interface{})

		chunkIndexValue, ok := properties["chunkIndex"].(float64)
		if !ok {
			return nil, fmt.Errorf("chunkIndex is not a float64 or is missing")
		}

		// Map the properties to your storemodels.Chunk struct
		singleChunkFromWeaviate := storemodels.Chunk{
			ChunkID:      uuid.MustParse(properties["chunkId"].(string)),
			UserID:       properties["userId"].(string),
			DocumentID:   uuid.MustParse(properties["documentId"].(string)),
			ChunkContent: properties["chunkContent"].(string),
			ChunkIndex:   int(chunkIndexValue),
		}

		chunksFromWeaviate = append(chunksFromWeaviate, singleChunkFromWeaviate)
	}
	return chunksFromWeaviate, nil
}

// func (w *WeaviateClientImpl) GetDocument(documentID string) (*Document, error) {
// 	objects, err := w.client.Data().ObjectsGetter().
// 		WithClassName("Documents").
// 		WithID(documentID).
// 		Do(context.Background())
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	// If no objects are returned, return an error
// 	if len(objects) == 0 {
// 		return nil, errors.New("no documents found")
// 	}
//
// 	// Combine chunks to form the complete document content
// 	var content string
// 	for _, obj := range objects {
// 		if obj.Properties == nil {
// 			return nil, errors.New("properties does not exist")
// 		}
//
// 		chunkValue, exists := obj.Properties.(map[string]interface{})["chunk"]
// 		if !exists || chunkValue == nil {
// 			return nil, errors.New("chunk does not exist")
// 		}
// 		content += chunkValue.(string)
// 	}
//
// 	// Assume the first object is the one you're looking for
// 	obj := objects[0]
//
// 	// Additional checks for each field before type assertion
// 	userID, ok := obj.Properties.(map[string]interface{})["userId"]
// 	if !ok || userID == nil {
// 		return nil, errors.New("userId does not exist")
// 	}
//
// 	documentName, ok := obj.Properties.(map[string]interface{})["documentName"]
// 	if !ok || documentName == nil {
// 		return nil, errors.New("documentName does not exist")
// 	}
//
// 	// Convert the object to a Document
// 	doc := &Document{
// 		UserID:       userID.(string),
// 		DocumentName: documentName.(string),
// 		Content:      content,
// 	}
//
// 	return doc, nil
// }

//	func (w *WeaviateClientImpl) UpdateDocumentContent(documentID, content string) error {
//		document := map[string]interface{}{
//			"content": content,
//		}
//
//		err := w.client.Data().Updater().
//			WithMerge().
//			WithID(documentID).
//			WithClassName("Documents").
//			WithProperties(document).
//			Do(context.Background())
//
//		return err
//	}
// func (w *WeaviateClientImpl) UpdateDocument(documentID, userID, name, content string) error {
// 	// First, delete all existing chunks for the document
// 	err := w.client.Data().Deleter().
// 		WithClassName("Documents").
// 		WithID(documentID).
// 		Do(context.Background())
// 	if err != nil {
// 		return err
// 	}
//
// 	// Now, use the UploadDocument function to add the new content
// 	err = w.UploadDocument(documentID, userID, name, content)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (w *WeaviateClientImpl) UpdateDocumentName(documentID, documentName string) error {
// 	document := map[string]interface{}{
// 		"documentName": documentName,
// 	}
//
// 	err := w.client.Data().Updater().
// 		WithMerge().
// 		WithID(documentID).
// 		WithClassName("Documents").
// 		WithProperties(document).
// 		Do(context.Background())
//
// 	return err
// }

// func (w *WeaviateClientImpl) SearchDocumentsByText(limit int, userID string, concepts []string) (*models.GraphQLResponse, error) {
// 	className := "Documents"
//
// 	documentName := graphql.Field{Name: "documentName"}
// 	content := graphql.Field{Name: "content"}
// 	_additional := graphql.Field{
// 		Name: "_additional", Fields: []graphql.Field{
// 			{Name: "certainty"}, // only supported if distance==cosine
// 			{Name: "distance"},  // always supported
// 		},
// 	}
//
// 	distance := float32(0.6)
// 	// moveAwayFrom := &graphql.MoveParameters{
// 	// 	Concepts: []string{"finance"},
// 	// 	Force:    0.45,
// 	// }
// 	// moveTo := &graphql.MoveParameters{
// 	// 	Concepts: []string{"haute couture"},
// 	// 	Force:    0.85,
// 	// }
// 	nearText := w.client.GraphQL().NearTextArgBuilder().
// 		WithConcepts(concepts).
// 		WithDistance(distance) // use WithCertainty(certainty) prior to v1.14
// 		// WithMoveTo(moveTo).
// 		// WithMoveAwayFrom(moveAwayFrom)
//
// 		// Creating the where filter
// 	whereFilter := filters.Where().
// 		WithPath([]string{"userId"}).
// 		WithOperator(filters.Equal).
// 		WithValueText(userID)
//
// 	ctx := context.Background()
//
// 	result, err := w.client.GraphQL().Get().
// 		WithClassName(className).
// 		WithFields(documentName, content, _additional).
// 		WithNearText(nearText).
// 		WithLimit(limit).
// 		WithWhere(whereFilter).
// 		Do(ctx)
//
// 	if err != nil {
// 		panic(err)
// 	}
// 	return result, nil
// }
