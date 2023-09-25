package documents

import (
	"bytes"
	"encoding/json"
	"lucidify-api/modules/clerkclient"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDocumentsUploadHandlerIntegration(t *testing.T) {
	testconfig := config.NewServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL
	db, err := store.NewStore(PostgresqlURL)
	// Setup the real environment
	clerkInstance, err := clerkclient.NewClerkClient(testconfig.ClerkSecretKey)

	if err != nil {
		t.Fatalf("Failed to create Clerk client: %v", err)
	}
	cfg := &config.ServerConfig{}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, db, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := testconfig.TestJWTSessionToken

	// Send a POST request to the server with the JWT token
	document := map[string]string{
		"document_name": "Test Document",
		"content":       "Test Content",
	}
	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/documents/upload", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Create the user with testconfig.user_id in the local database
	// the user id registered by the jwt token must exist in the local database
	user := store.User{
		UserID:           testconfig.TestUserID,
		ExternalID:       "TestCreateUserInUsersTableExternalID",
		Username:         "TestCreateUserInUsersTableUsername",
		PasswordEnabled:  true,
		Email:            "TestCreateUserInUsersTable@example.com",
		FirstName:        "TestCreateUserInUsersTableCreateTest",
		LastName:         "TestCreateUserInUsersTableUser",
		ImageURL:         "https://TestCreateUserInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://TestCreateUserInUsersTable.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err = db.CreateUserInUsersTable(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Check if the user exists
	err = db.CheckIfUserInUsersTable(user.UserID, 3)
	if err != nil {
		t.Fatalf("User not found after creation: %v", err)
	}

	db.UploadDocument(testconfig.TestUserID, "Test Document", "Test Content")
	documentFromDb, err := db.GetDocument(testconfig.TestUserID, "Test Document")
	if err != nil {
		t.Fatalf("Failed to get document: %v", err)
	}

	documentFromDbContent := documentFromDb.Content
	if documentFromDbContent != "Test Content" {
		t.Fatalf("Expected document content %s, got %s", "Test Content", documentFromDbContent)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		db.DeleteUserInUsersTable(UserID)
		db.DeleteDocument(UserID, "Test Document")
	})
}

func TestDocumentsUploadHandlerUnauthorizedIntegration(t *testing.T) {
	testconfig := config.NewServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL
	db, err := store.NewStore(PostgresqlURL)
	// Setup the real environment
	clerkInstance, err := clerkclient.NewClerkClient(testconfig.ClerkSecretKey)

	if err != nil {
		t.Fatalf("Failed to create Clerk client: %v", err)
	}
	cfg := &config.ServerConfig{}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, db, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := testconfig.TestJWTSessionToken + "invalid"

	// Send a POST request to the server with the JWT token
	document := map[string]string{
		"document_name": "Test Document",
		"content":       "Test Content",
	}
	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/documents/upload", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected status code not OK, got %d", resp.StatusCode)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		db.DeleteUserInUsersTable(UserID)
		db.DeleteDocument(UserID, "Test Document")
	})
}
