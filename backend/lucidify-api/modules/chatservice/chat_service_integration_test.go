// //go:build integration
// // +build integration
package chatservice

import (
	"log"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"
	"testing"
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

	// Create instance of ChatService
	chatService := NewChatService(postgresqlDB, weaviateDB)

	return chatService
}

func TestProcessCurrentThreadAndReturnSystemPromptIntegration(t *testing.T) {
	// Setup ChatService
	chatService := setupTestChatService()

	// Here you would set up any necessary preconditions in your database,
	// like creating chat threads that you're going to process.

	// You might want to use a real thread from your database, or create a new one just for testing.
	// Assuming here that you've set up a "current thread" in some way.

	// Call the method to test.
	response, err := chatService.ProcessCurrentThreadAndReturnSystemPrompt()

	// Assert that the expected response is returned.
	// This will depend on what exactly your method does and what your expected outcomes are.
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
