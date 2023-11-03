package weaviateclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/server/config"

	"github.com/google/uuid"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type WeaviateClient interface {
	GetWeaviateClient() *weaviate.Client
	UploadChunk(storemodels.Chunk) error
	UploadChunks([]storemodels.Chunk) error
	DeleteChunk(chunkID uuid.UUID) error
	DeleteChunks([]storemodels.Chunk) error
	DeleteAllChunksByUserID(userID string) error
	GetChunks(chunksFromPostgresql []storemodels.Chunk) ([]storemodels.Chunk, error)
	SearchDocumentsByText(limit int, userID string, concepts []string) ([]storemodels.ChunkFromVectorSearch, error)
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
		Host:   "weaviate:8080",
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

func NewWeaviateClientTest() (WeaviateClient, error) {
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
		WithID(chunk.ChunkID.String()).
		WithClassName("Documents").
		WithProperties(chunkData).
		Do(context.Background())

	if err != nil {
		return fmt.Errorf("failed to upload chunk: %w", err)
	}

	return nil
}

func (w *WeaviateClientImpl) DeleteChunk(chunkID uuid.UUID) error {
	err := w.client.Data().Deleter().
		WithClassName("Documents").
		WithID(chunkID.String()).
		Do(context.Background())

	return err
}

func (w *WeaviateClientImpl) DeleteAllChunksByUserID(userID string) error {
	if w.client == nil {
		log.Println("Client is nil in deleteDocumentsByUserID")
		return fmt.Errorf("client is nil")
	}

	// Define the where filter to match all documents with the given userId
	whereFilter := filters.Where().
		WithPath([]string{"userId"}).
		WithOperator(filters.Equal).
		WithValueText(userID)

	// Perform the batch delete operation
	response, err := w.client.Batch().ObjectsBatchDeleter().
		WithClassName("Documents").
		WithOutput("verbose").
		WithWhere(whereFilter).
		Do(context.Background())

	if err != nil {
		// Handle error
		// panic(err)
	}

	// Process the response
	fmt.Printf("Delete response: %+v\n", *response)
	return nil
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

func (w *WeaviateClientImpl) SearchDocumentsByText(limit int, userID string, concepts []string) ([]storemodels.ChunkFromVectorSearch, error) {
	className := "Documents"

	documentId := graphql.Field{Name: "documentId"}
	chunkId := graphql.Field{Name: "chunkId"}
	chunkContent := graphql.Field{Name: "chunkContent"}
	chunkIndex := graphql.Field{Name: "chunkIndex"}
	_additional := graphql.Field{
		Name: "_additional", Fields: []graphql.Field{
			{Name: "certainty"}, // only supported if distance==cosine
			{Name: "distance"},  // always supported
		},
	}

	distance := float32(0.6)
	// moveAwayFrom := &graphql.MoveParameters{
	// 	Concepts: []string{"finance"},
	// 	Force:    0.45,
	// }
	// moveTo := &graphql.MoveParameters{
	// 	Concepts: []string{"haute couture"},
	// 	Force:    0.85,
	// }
	nearText := w.client.GraphQL().NearTextArgBuilder().
		WithConcepts(concepts).
		WithDistance(distance) // use WithCertainty(certainty) prior to v1.14
	// WithMoveTo(moveTo).
	// WithMoveAwayFrom(moveAwayFrom)

	// Creating the where filter
	whereFilter := filters.Where().
		WithPath([]string{"userId"}).
		WithOperator(filters.Equal).
		WithValueText(userID)

	ctx := context.Background()

	result, err := w.client.GraphQL().Get().
		WithClassName(className).
		WithFields(documentId, chunkId, chunkContent, chunkIndex, _additional).
		WithNearText(nearText).
		WithLimit(limit).
		WithWhere(whereFilter).
		Do(ctx)

	if err != nil {
		panic(err)
	}

	var chunks []storemodels.ChunkFromVectorSearch

	if result != nil && result.Data != nil {
		getData, ok := result.Data["Get"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected format for 'Get' data")
		}

		unprocessedChunks, ok := getData["Documents"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected format for 'Documents' data")
		}

		for _, chunk := range unprocessedChunks {
			docMap, ok := chunk.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("unexpected format for chunk data")
			}

			// documentName := docMap["documentName"].(string)
			documentId := docMap["documentId"].(string)
			chunkId := docMap["chunkId"].(string)
			chunkContent := docMap["chunkContent"].(string)
			chunkIndex := docMap["chunkIndex"].(float64)
			additional := docMap["_additional"].(map[string]interface{})
			certainty := additional["certainty"].(float64)
			distance := additional["distance"].(float64)

			chunkFromVectorSearch := storemodels.ChunkFromVectorSearch{
				ChunkID:      uuid.MustParse(chunkId),
				UserID:       userID,
				DocumentID:   uuid.MustParse(documentId),
				ChunkContent: chunkContent,
				ChunkIndex:   int(chunkIndex),
				Certainty:    certainty,
				Distance:     distance,
			}
			chunks = append(chunks, chunkFromVectorSearch)
		}
	}

	return chunks, nil
}
