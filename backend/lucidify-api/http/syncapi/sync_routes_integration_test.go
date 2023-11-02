// //go:build integration
// // +build integration
package syncapi

import (
	"bytes"
	"encoding/json"
	"io"
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

	// Authenticated request
	req, _ := http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=conversationHistory",
		bytes.NewBuffer(body))

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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"message":"Data set successfully for key: conversationHistory"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"message":"Data set successfully for key: conversationHistory"}`,
			string(respBody))
	}

	// Autheticated Request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=conversationHistory",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"data":"\"conversationHistory\"","message":"Data fetched successfully"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"data":"\"conversationHistory\"","message":"Data fetched successfully"}`, string(respBody))
	}

	// Testing the on conflict (postgres) functionality to update chat
	// Send a POST request to the server with the JWT token
	body, _ = json.Marshal("conversationHistory but updated!")

	// Authenticated request
	req, _ = http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=conversationHistory",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"message":"Data set successfully for key: conversationHistory"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"message":"Data set successfully for key: conversationHistory"}`,
			string(respBody))
	}

	// Autheticated Request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=conversationHistory",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"data":"\"conversationHistory but updated!\"","message":"Data fetched successfully"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"data":"\"conversationHistory but updated!\"","message":"Data fetched successfully"}`, string(respBody))
	}

	// Unauthenticated POST request
	req, _ = http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=conversationHistory",
		bytes.NewBuffer(body))

	req.Header.Set("Authorization", "Bearer "+jwtToken+"invalid")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	// This is counter intuitive. Not sure why clerk doesn't return a 401.
	if string(respBody) != `couldn't find cookie __session` {
		t.Errorf("Expected response body %s, got %s", `couldn't find cookie __session`, string(respBody))
	}
	// Unauthenticated GET request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=conversationHistory",
		bytes.NewBuffer(body))

	req.Header.Set("Authorization", "Bearer "+jwtToken+"invalid")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	// This is counter intuitive. Not sure why clerk doesn't return a 401.
	if string(respBody) != `couldn't find cookie __session` {
		t.Errorf("Expected response body %s, got %s", `couldn't find cookie __session`, string(respBody))
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
	body, _ := json.Marshal("someFoldersData")

	// Authenticated request
	req, _ := http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=folders",
		bytes.NewBuffer(body))

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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"message":"Data set successfully for key: folders"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"message":"Data set successfully for key: folders"}`,
			string(respBody))
	}

	// Autheticated Request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=folders",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"data":"\"someFoldersData\"","message":"Data fetched successfully"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"data":"\"someFoldersData\"","message":"Data fetched successfully"}`, string(respBody))
	}

	// Testing the on conflict (postgres) functionality to update chat
	// Send a POST request to the server with the JWT token
	body, _ = json.Marshal("someFoldersData but updated!")

	// Authenticated request
	req, _ = http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=folders",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"message":"Data set successfully for key: folders"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"message":"Data set successfully for key: folders"}`,
			string(respBody))
	}

	// Autheticated Request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=folders",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"data":"\"someFoldersData but updated!\"","message":"Data fetched successfully"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"data":"\"someFoldersData but updated!\"","message":"Data fetched successfully"}`, string(respBody))
	}

	// Unauthenticated POST request
	req, _ = http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=folders",
		bytes.NewBuffer(body))

	req.Header.Set("Authorization", "Bearer "+jwtToken+"invalid")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	// This is counter intuitive. Not sure why clerk doesn't return a 401.
	if string(respBody) != `couldn't find cookie __session` {
		t.Errorf("Expected response body %s, got %s", `couldn't find cookie __session`, string(respBody))
	}
	// Unauthenticated GET request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=folders",
		bytes.NewBuffer(body))

	req.Header.Set("Authorization", "Bearer "+jwtToken+"invalid")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	// This is counter intuitive. Not sure why clerk doesn't return a 401.
	if string(respBody) != `couldn't find cookie __session` {
		t.Errorf("Expected response body %s, got %s", `couldn't find cookie __session`, string(respBody))
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
	body, _ := json.Marshal("somePromptsData")

	// Authenticated request
	req, _ := http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=prompts",
		bytes.NewBuffer(body))

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
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"message":"Data set successfully for key: prompts"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"message":"Data set successfully for key: prompts"}`,
			string(respBody))
	}

	// Autheticated Request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=prompts",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"data":"\"somePromptsData\"","message":"Data fetched successfully"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"data":"\"somePromptsData\"","message":"Data fetched successfully"}`,
			string(respBody))
	}

	// Testing the on conflict (postgres) functionality to update chat
	// Send a POST request to the server with the JWT token
	body, _ = json.Marshal("somePromptsData but updated!")

	// Authenticated request
	req, _ = http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=prompts",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"message":"Data set successfully for key: prompts"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"message":"Data set successfully for key: prompts"}`,
			string(respBody))
	}

	// Autheticated Request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=prompts",
		bytes.NewBuffer(body))

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
	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	if string(respBody) !=
		`{"success":true,"data":"\"somePromptsData but updated!\"","message":"Data fetched successfully"}` {
		t.Errorf("Expected response body %s, got %s",
			`{"success":true,"data":"\"somePromptsData but updated!\"","message":"Data fetched successfully"}`,
			string(respBody))
	}

	// Unauthenticated POST request
	req, _ = http.NewRequest(
		http.MethodPost,
		server.URL+"/api/sync/localstorage/?key=prompts",
		bytes.NewBuffer(body))

	req.Header.Set("Authorization", "Bearer "+jwtToken+"invalid")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	// This is counter intuitive. Not sure why clerk doesn't return a 401.
	if string(respBody) != `couldn't find cookie __session` {
		t.Errorf("Expected response body %s, got %s", `couldn't find cookie __session`, string(respBody))
	}
	// Unauthenticated GET request
	req, _ = http.NewRequest(
		http.MethodGet,
		server.URL+"/api/sync/localstorage/?key=prompts",
		bytes.NewBuffer(body))

	req.Header.Set("Authorization", "Bearer "+jwtToken+"invalid")
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Read the response body
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	// This is counter intuitive. Not sure why clerk doesn't return a 401.
	if string(respBody) != `couldn't find cookie __session` {
		t.Errorf("Expected response body %s, got %s",
			`couldn't find cookie __session`, string(respBody))
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		// postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}
