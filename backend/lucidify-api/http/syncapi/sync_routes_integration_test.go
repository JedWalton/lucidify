// //go:build integration
// // +build integration
package syncapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/server/config"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/syncservice"
	"lucidify-api/service/userservice"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func createTestUserInDb() error {
	testconfig := config.NewServerConfig()
	db, err := postgresqlclient.NewPostgreSQL()

	// the user id registered by the jwt token must exist in the local database
	user := storemodels.User{
		UserID:           testconfig.TestUserID,
		ExternalID:       "TestCreateUserInUsersTableExternalIDDocuments",
		Username:         "TestCreateUserInUsersTableUsernameDocuments",
		PasswordEnabled:  true,
		Email:            "TestCreateUserInUsersTableDocuments@example.com",
		FirstName:        "TestCreateUserInUsersTableCreateTest",
		LastName:         "TestCreateUserInUsersTableUser",
		ImageURL:         "https://TestCreateUserInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://TestCreateUserInUsersTable.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	userService, err := userservice.NewUserService()
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

	err = db.CreateUserInUsersTable(user)
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

	err = createTestUserInDb()
	if err != nil {
		t.Fatalf("Failed to create test user in db: %v", err)
	}

	return &TestSetup{
		Config:        cfg,
		PostgresqlDB:  postgresqlDB,
		ClerkInstance: clerkInstance,
	}
}

func TestConversationHistoryIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	syncService, err := syncservice.NewSyncService()
	if err != nil {
		t.Fatalf("Failed to create SyncService: %v", err)
	}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, clerkInstance, syncService)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken

	// Send a POST request to the server with the JWT token
	body, _ := json.Marshal("conversationHistory")
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/sync/localstorage/?key=conversationHistory", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	// Read the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) != `{"success":true,"message":"Data set successfully for key: conversationHistory"}` {
		t.Errorf("Expected response body %s, got %s", `{"success":true,"message":"Data set successfully for key: conversationHistory"}`, string(respBody))
	}
	// Print the response body
	// t.Fatalf("Response Body: %s", respBody)

	body, _ = json.Marshal("conversationHistory")
	req, _ = http.NewRequest(http.MethodGet, server.URL+"/api/sync/localstorage/?key=conversationHistory", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		// postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

func TestPromptsIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	syncService, err := syncservice.NewSyncService()
	if err != nil {
		t.Fatalf("Failed to create SyncService: %v", err)
	}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, clerkInstance, syncService)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken

	// Send a POST request to the server with the JWT token
	body, _ := json.Marshal("prompts")
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/sync/localstorage/?key=prompts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Send a GET request to the server with the JWT token
	body, _ = json.Marshal("prompts")
	req, _ = http.NewRequest(http.MethodGet, server.URL+"/api/sync/localstorage/?key=prompts", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		// postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

func TestFoldersIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	syncService, err := syncservice.NewSyncService()
	if err != nil {
		t.Fatalf("Failed to create SyncService: %v", err)
	}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, clerkInstance, syncService)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken

	// Send a POST request to the server with the JWT token
	body, _ := json.Marshal("folders")
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/api/sync/localstorage/?key=folders", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	body, _ = json.Marshal("folders")
	req, _ = http.NewRequest(http.MethodGet, server.URL+"/api/sync/localstorage/?key=folders", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		// postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

// func TestClearConversationsFromServerIntegration(t *testing.T) {
// 	setup := setupTestServer(t)
// 	cfg := setup.Config
// 	syncService := setup.SyncService
// 	clerkInstance := setup.ClerkInstance
//
// 	// Create a test server
// 	mux := http.NewServeMux()
// 	SetupRoutes(cfg, mux, clerkInstance, syncService)
// 	server := httptest.NewServer(mux)
// 	defer server.Close()
//
// 	// Obtain a JWT token from Clerk
// 	jwtToken := cfg.TestJWTSessionToken
//
// 	// Send a DELETE request to the server with the JWT token
// 	req, err := http.NewRequest(http.MethodDelete, server.URL+"/api/sync/localstorage/?key=clearConversations", bytes.NewBuffer(nil))
// 	if err != nil {
// 		t.Errorf("Failed to create request: %v", err)
// 	}
// 	req.Header.Set("Authorization", "Bearer "+jwtToken)
//
// 	// Add the "__session" cookie to the request
// 	// sessionCookieValue := jwtToken
// 	// req.AddCookie(&http.Cookie{Name: "__session", Value: sessionCookieValue})
//
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Errorf("Failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()
//
// 	// Check the response
// 	if resp.StatusCode != http.StatusOK {
// 		bodyBytes, _ := io.ReadAll(resp.Body)
// 		t.Errorf("Expected status code %d, got %d. Response body: %s", http.StatusOK, resp.StatusCode, bodyBytes)
// 	}
//
// 	// Check if the conversations are actually cleared
// 	// TODO: Implement a function to verify that the conversations for the user have been cleared.
// 	// For example, you can use a DB function like `GetConversations(userID)` and ensure the list is empty.
// }
