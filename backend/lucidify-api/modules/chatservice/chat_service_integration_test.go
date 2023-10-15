// //go:build integration
// // +build integration
package chatservice

import (
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/store"
	"lucidify-api/modules/store/weaviateclient"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func createTestUserInDb() string {
	db, err := postgresqlclient.NewPostgreSQL()

	// the user id registered by the jwt token must exist in the local database
	user := postgresqlclient.User{
		UserID:           "TestChatServiceIntegrationTestUUID",
		ExternalID:       "TestChatServiceIntegrationTestExternalID",
		Username:         "TestChatServiceIntegrationTestUsername",
		PasswordEnabled:  true,
		Email:            "TestChatServiceIntTest@gmail.com",
		FirstName:        "TestChatServiceIntegrationTestFirstName",
		LastName:         "TestChatServiceIntegrationTestLastName",
		ImageURL:         "https://TestChatServiceIntegrationTestURL.com/image.jpg",
		ProfileImageURL:  "https://TestChatServiceTestProfileURL.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	db.DeleteUserInUsersTable(user.UserID)
	err = db.CheckUserDeletedInUsersTable(user.UserID, 3)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	err = db.CreateUserInUsersTable(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Check if the user exists
	err = db.CheckIfUserInUsersTable(user.UserID, 3)
	if err != nil {
		log.Fatalf("User not found after creation: %v", err)
	}
	return user.UserID
}

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

	documentService := store.NewDocumentService(postgresqlDB, weaviateDB)

	// Create instance of ChatService
	chatService := NewChatService(postgresqlDB, weaviateDB, openaiClient, documentService)

	createTestUserInDb()

	documentService.UploadDocument("TestChatServiceIntegrationTestUUID", "Erica cat lore", "Erica likes cats, now we are mentioning cats erica is paying attention.")
	documentService.UploadDocument("TestChatServiceIntegrationTestUUID", "Erica cat 2 lore", "Erica likes Yuki (her new cat), now we are mentioning Yuki, erica is paying attention.")

	return chatService
}

// func TestChatCompletion(t *testing.T) {
// 	chatService := setupTestChatService()
//
// 	response, err := chatService.ChatCompletion("TestChatServiceIntegrationTestUUID")
//
// 	if err != nil {
// 		t.Errorf("Error was not expected while processing current thread: %v", err)
// 	}
//
// 	expectedResponse := "PLACEHOLDER RESPONSE" // Adjust "EXPECTED RESPONSE" to match what you're actually expecting.
// 	if response != expectedResponse {
// 		t.Errorf("Unexpected response: got %v want %v", response, expectedResponse)
// 	}
//
// 	// Optionally, you might want to query your databases here to assert that the expected
// 	// updates have been made as a result of calling the method.
//
// 	// Cleanup after test
// 	// Here you would clean up your database from any records you created for your test.
// }

func TestGetAnswerFromFiles(t *testing.T) {
	chatService := setupTestChatService()

	response, err := chatService.GetAnswerFromFiles("Who is erica and does she have pets?", "TestChatServiceIntegrationTestUUID")

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
