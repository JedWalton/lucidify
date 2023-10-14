// //go:build integration
// // +build integration
package chatservice

import (
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func setupTestChatService() ChatService {
	// Initialize PostgreSQL for tests
	postgresqlDB, err := postgresqlclient.NewPostgreSQL() // Adjust this to match your actual constructor
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	// Initialize Weaviate for tests
	weaviateDB, err := weaviateclient.NewWeaviateClientTest() // Adjust this to match your actual constructor
	if err != nil {
		log.Fatalf("Failed to create Weaviate client: %v", err)
	}

	cfg := config.NewServerConfig()
	openaiClient := openai.NewClient(cfg.OPENAI_API_KEY)

	// Create instance of ChatService
	chatService := NewChatService(postgresqlDB, weaviateDB, openaiClient)

	return chatService
}

func TestChatCompletion(t *testing.T) {
	chatService := setupTestChatService()

	response, err := chatService.ChatCompletion()

	if err != nil {
		t.Errorf("Error was not expected while processing current thread: %v", err)
	}

	expectedResponse := "PLACEHOLDER RESPONSE" // Adjust "EXPECTED RESPONSE" to match what you're actually expecting.
	if response != expectedResponse {
		t.Errorf("Unexpected response: got %v want %v", response, expectedResponse)
	}

	// Optionally, you might want to query your databases here to assert that the expected
	// updates have been made as a result of calling the method.

	// Cleanup after test
	// Here you would clean up your database from any records you created for your test.
}
