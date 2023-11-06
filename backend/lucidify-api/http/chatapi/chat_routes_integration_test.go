// //go:build integration
// // +build integration
package chatapi

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/server/config"
	"lucidify-api/service/chatservice"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
	"lucidify-api/service/userservice"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/sashabaranov/go-openai"
)

func createTestUserInDb(cfg *config.ServerConfig, db *postgresqlclient.PostgreSQL) error {
	// the user id registered by the jwt token must exist in the local database
	user := storemodels.User{
		UserID:           cfg.TestUserID,
		ExternalID:       "TestChatAPIUserInUsersTableExternalIDDocuments",
		Username:         "TestChatAPIUsersTableUsernameDocuments",
		PasswordEnabled:  true,
		Email:            "TestChatAPIUserUsersTableDocuments@example.com",
		FirstName:        "TestUsersTableCreateTest",
		LastName:         "TestUsersTableUser",
		ImageURL:         "https://TestInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://TestInUsersTable.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	weaviate, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		log.Fatalf("Failed to create WeaviateClient: %v", err)
	}
	userService, err := userservice.NewUserService(db, weaviate)
	if err != nil {
		log.Fatalf("Failed to create UserService: %v", err)
	}

	err = userService.DeleteUser(user.UserID)
	if err != nil {
		log.Fatalf("Failed to delete user: %v", err)
	}
	if !userService.HasUserBeenDeleted(user.UserID, 3) {
		log.Fatalf("Failed to delete user: %v", err)
	}

	err = userService.CreateUser(user)
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	// Check if the user exists
	_, err = userService.GetUserWithRetries(user.UserID, 3)
	if err != nil {
		log.Fatalf("User not found after creation: %v", err)
	}

	return nil
}

type TestSetup struct {
	Config        *config.ServerConfig
	PostgresqlDB  *postgresqlclient.PostgreSQL
	ClerkInstance clerk.Client
	Weaviate      weaviateclient.WeaviateClient
	DocService    documentservice.DocumentService
}

func SetupTestEnvironment(t *testing.T) *TestSetup {
	cfg := config.NewServerConfig()

	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		t.Fatalf("Failed to create test postgresqlclient: %v", err)
	}

	clerkInstance, err := clerkservice.NewClerkClient()
	if err != nil {
		t.Fatalf("Failed to create Clerk client: %v", err)
	}

	weaviate, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		t.Fatalf("Failed to create WeaviateClient: %v", err)
	}

	docService := documentservice.NewDocumentService(postgresqlDB, weaviate)

	err = createTestUserInDb(cfg, postgresqlDB)
	if err != nil {
		t.Fatalf("Failed to create test user in db: %v", err)
	}

	return &TestSetup{
		Config:        cfg,
		PostgresqlDB:  postgresqlDB,
		ClerkInstance: clerkInstance,
		Weaviate:      weaviate,
		DocService:    docService,
	}
}

func TestChatHandlerIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	clerkInstance := setup.ClerkInstance
	openaiClient := openai.NewClient(cfg.OPENAI_API_KEY)
	documentService := setup.DocService
	chatVectorService := chatservice.NewChatVectorService(setup.Weaviate, openaiClient, documentService)

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, chatVectorService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()
}

//
// 	// Obtain a JWT token from Clerk
// 	// jwtToken := cfg.TestJWTSessionToken
//
// 	// Construct a message
// 	messages := []Message{
// 		{Role: RoleUser, Content: "Hello, how can I help you?"},
// 	}
//
// 	// Send a POST request to the server with the JWT token and message
// 	// body, _ := json.Marshal(map[string][]Message{"messages": messages})
// 	_, _ = json.Marshal(map[string][]Message{"messages": messages})
//
// 	// // Authenticated request
// 	// req, _ := http.NewRequest(
// 	// 	http.MethodPost,
// 	// 	server.URL+"/api/chat/vector-search",
// 	// 	bytes.NewBuffer(body))
// 	//
// 	// req.Header.Set("Authorization", "Bearer "+jwtToken)
// 	// client := &http.Client{}
// 	// resp, err := client.Do(req)
// 	// if err != nil {
// 	// 	t.Errorf("Failed to send request: %v", err)
// 	// }
// 	// defer resp.Body.Close()
// 	//
// 	// // Check the response
// 	// if resp.StatusCode != http.StatusOK {
// 	// 	t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
// 	// }
// 	// // Read the response body
// 	// respBody, err := io.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	t.Errorf("Failed to read response body: %v", err)
// 	// }
// 	//
// 	// var serverResp ServerResponse
// 	// err = json.Unmarshal(respBody, &serverResp)
// 	// if err != nil {
// 	// 	t.Errorf("Failed to unmarshal response body: %v", err)
// 	// }
// 	//
// 	// if !serverResp.Success || serverResp.Message == "" || serverResp.Data == nil {
// 	// 	t.Errorf("Expected successful response with data and message, got %+v", serverResp)
// 	// }
//
// 	// Additional tests can be performed here, like checking the specific content of the systemPrompt
// 	// and ensuring it matches expected values based on the input "messages".
// 	//
// 	// // Unauthenticated POST request
// 	// req, _ = http.NewRequest(
// 	// 	http.MethodPost,
// 	// 	server.URL+"/api/chat/vector-search",
// 	// 	bytes.NewBuffer(body))
// 	//
// 	// req.Header.Set("Authorization", "Bearer "+jwtToken+"invalid")
// 	// client = &http.Client{}
// 	// resp, err = client.Do(req)
// 	// if err != nil {
// 	// 	t.Errorf("Failed to send request: %v", err)
// 	// }
// 	// defer resp.Body.Close()
// 	//
// 	// // Check the response
// 	// if resp.StatusCode != http.StatusUnauthorized {
// 	// 	t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
// 	// }
// 	//
// 	// // Read the response body
// 	// respBody, err = io.ReadAll(resp.Body)
// 	// if err != nil {
// 	// 	t.Errorf("Failed to read response body: %v", err)
// 	// }
// 	//
// 	// if string(respBody) != "Unauthorized" {
// 	// 	t.Errorf("Expected response body %s, got %s", "Unauthorized", string(respBody))
// 	// }
// 	//
// 	// // Cleanup if necessary
// 	// t.Cleanup(func() {
// 	// 	testconfig := config.NewServerConfig()
// 	// 	UserID := testconfig.TestUserID
// 	// 	postgresqlDB := setup.PostgresqlDB
// 	// 	postgresqlDB.DeleteUserInUsersTable(UserID)
// 	// })
// }

// SetupTestEnvironment is assumed to be the same as in the previous example, which prepares
// the database and other services for integration testing.
