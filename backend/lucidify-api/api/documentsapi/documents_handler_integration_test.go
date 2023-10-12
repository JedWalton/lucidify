// // go:build integration
// // +build integration
package documentsapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"lucidify-api/modules/clerkclient"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/store"
	"lucidify-api/modules/store/storemodels"
	"lucidify-api/modules/store/weaviateclient"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func createTestUserInDb() error {
	testconfig := config.NewServerConfig()
	db, err := postgresqlclient.NewPostgreSQL()

	// the user id registered by the jwt token must exist in the local database
	user := postgresqlclient.User{
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
		return err
	}

	return nil
}
func createASecondTestUserInDb() string {
	db, err := postgresqlclient.NewPostgreSQL()

	user := postgresqlclient.User{
		UserID:           "userid_testuserid2",
		ExternalID:       "TestCreateSecondUserInUsersTableExternalID",
		Username:         "TestCreateSecondUserInUsersTableUsername",
		PasswordEnabled:  true,
		Email:            "TestCreateSecondUserInUsersTable@example.com",
		FirstName:        "TestCreateSecondUserInUsersTableCreateTest",
		LastName:         "TestCreateSecondUserInUsersTableUser",
		ImageURL:         "https://TestCreateSecondUserInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://TestCreateSecondUserInUsersTable.com/profile.jpg",
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

type TestSetup struct {
	Config          *config.ServerConfig
	PostgresqlDB    *postgresqlclient.PostgreSQL
	ClerkInstance   clerk.Client
	WeaviateDB      weaviateclient.WeaviateClient
	DocumentService store.DocumentService
}

func SetupTestEnvironment(t *testing.T) *TestSetup {
	cfg := config.NewServerConfig()

	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		t.Fatalf("Failed to create test postgresqlclient: %v", err)
	}

	clerkInstance, err := clerkclient.NewClerkClient(cfg.ClerkSecretKey)
	if err != nil {
		t.Fatalf("Failed to create Clerk client: %v", err)
	}

	weaviateDB, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		t.Fatalf("Failed to create Weaviate client: %v", err)
	}

	err = createTestUserInDb()
	if err != nil {
		t.Fatalf("Failed to create test user in db: %v", err)
	}

	documentService := store.NewDocumentService(postgresqlDB, weaviateDB)

	return &TestSetup{
		Config:          cfg,
		PostgresqlDB:    postgresqlDB,
		ClerkInstance:   clerkInstance,
		WeaviateDB:      weaviateDB,
		DocumentService: documentService,
	}
}

func TestDocumentsUploadHandlerIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken

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
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	documentFromDb, err := postgresqlDB.GetDocument(cfg.TestUserID, "Test Document")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}

	documentFromDbContent := documentFromDb.Content
	if documentFromDbContent != "Test Content" {
		t.Errorf("Expected document content %s, got %s", "Test Content", documentFromDbContent)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

func TestDocumentsUploadHandlerUnauthorizedIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken + "invalid"

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
		t.Errorf("Failed to send request: %v", err)
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
		postgresqlDB.DeleteUserInUsersTable(UserID)
		postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

func TestDocumentsGetDocumentHandlerIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	jwtToken := cfg.TestJWTSessionToken

	document := map[string]string{
		"document_name": "Test Document",
		"content":       "Test Content",
	}

	postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")

	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/documents/getdocument", bytes.NewBuffer(body))
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

	// Please implement the rest of this integration test to check it returns the correct document
	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body into a Document object
	var respDocument storemodels.Document
	err = json.Unmarshal(respBody, &respDocument)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	// Check if the returned document is correct
	if respDocument.DocumentName != document["document_name"] || respDocument.Content != document["content"] {
		t.Errorf("Returned document does not match the expected document")
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

func TestDocumentsGetDocumentHandlerUnauthorizedIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	jwtToken := cfg.TestJWTSessionToken + "invalid"

	document := map[string]string{
		"document_name": "Test Document",
		"content":       "Test Content",
	}

	postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")

	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/documents/getdocument", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code Bad Request, 400. Got: %v", resp.StatusCode)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

func TestDocumentsGetAllDocumentsHandlerIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	jwtToken := cfg.TestJWTSessionToken

	postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")
	postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document 2", "Test Content 2")
	postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document 3", "Test Content 3")

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/documents/get_all_documents", nil)
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

	// Unmarshal the response body into a slice of Document objects
	var respDocuments []storemodels.Document
	err = json.Unmarshal(respBody, &respDocuments)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	// Check if the returned documents are correct
	if len(respDocuments) != 3 {
		t.Errorf("Expected 3 documents, got %d", len(respDocuments))
	}

	expectedDocs := []string{"Test Document", "Test Document 2", "Test Document 3"}
	for i, doc := range respDocuments {
		if doc.DocumentName != expectedDocs[i] {
			t.Errorf("Expected document name %s, got %s", expectedDocs[i], doc.DocumentName)
		}
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
	})
}

func TestDocumentsGetAllDocumentsHandlerUnauthorizedIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	jwtToken := cfg.TestJWTSessionToken + " invalid"

	postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")
	postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document 2", "Test Content 2")
	postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document 3", "Test Content 3")

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/documents/get_all_documents", nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
	})
}

func TestDocumentsGetAllDocumentsHandlerUnauthenticatedOtherUserIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	UserID2 := createASecondTestUserInDb()

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	jwtToken := cfg.TestJWTSessionToken

	postgresqlDB.UploadDocument(UserID2, "Test Document", "Test Content")
	postgresqlDB.UploadDocument(UserID2, "Test Document 2", "Test Content 2")
	postgresqlDB.UploadDocument(UserID2, "Test Document 3", "Test Content 3")

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/documents/get_all_documents", nil)
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to read response body: %v", err)
	}

	// Unmarshal the response body into a slice of Document objects
	var respDocuments []storemodels.Document
	err = json.Unmarshal(respBody, &respDocuments)
	if err != nil {
		t.Errorf("Failed to unmarshal response body: %v", err)
	}

	// Check if the returned documents are correct
	if len(respDocuments) != 0 {
		t.Errorf("Expected 0 documents, got %d", len(respDocuments))
	}

	expectedDocs := []string{"Test Document", "Test Document 2", "Test Document 3"}
	for i, doc := range respDocuments {
		if doc.DocumentName == expectedDocs[i] {
			t.Errorf("Expected to not mach document name %s, got %s", expectedDocs[i], doc.DocumentName)
		}
	}
	// Cleanup the database
	t.Cleanup(func() {
		postgresqlDB.DeleteUserInUsersTable(cfg.TestUserID)
		postgresqlDB.DeleteUserInUsersTable(UserID2)
		postgresqlDB.DeleteDocument(UserID2, "Test Document")
		postgresqlDB.DeleteDocument(UserID2, "Test Document 2")
		postgresqlDB.DeleteDocument(UserID2, "Test Document 3")
	})
}

func TestDocumentsDeleteDocumentHandlerIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	createTestUserInDb()

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken

	document, err := postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")
	if err != nil {
		t.Errorf("Failed to upload document: %v", err)
	}

	// Send a POST request to the server with the JWT token
	documentBodyRaw := map[string]string{
		"documentID": document.DocumentUUID.String(),
	}
	body, _ := json.Marshal(documentBodyRaw)
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/documents/deletedocument", bytes.NewBuffer(body))
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

	_, err = postgresqlDB.GetDocument(cfg.TestUserID, "Test Document")
	if err == nil {
		t.Errorf("Should have failed to get document: %v", err)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

func TestDocumentsDeleteDocumentHandlerNotMyDocumentIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	createASecondTestUserInDb()

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken

	document, err := postgresqlDB.UploadDocument("userid_testuserid2", "Test Document", "Test Content")
	if err != nil {
		t.Errorf("Failed to upload document: %v", err)
	}

	// Send a POST request to the server with the JWT token
	documentBodyRaw := map[string]string{
		"documentID": document.DocumentUUID.String(),
	}
	body, _ := json.Marshal(documentBodyRaw)
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/documents/deletedocument", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected status code to not be ok. %d, got %d", http.StatusOK, resp.StatusCode)
	}

	_, err = postgresqlDB.GetDocument("userid_testuserid2", "Test Document")
	if err != nil {
		t.Errorf("Should have not failed to get document as document should have not been deleted: %v", err)
	}

	// Cleanup the database
	t.Cleanup(func() {
		UserID := cfg.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		SecondUserID := "userid_testuserid2"
		postgresqlDB.DeleteUserInUsersTable(SecondUserID)
	})
}

func TestDocumentsDeleteDocumentHandlerUnauthenticatedIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	err := createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	// Obtain a JWT token from Clerk
	jwtToken := cfg.TestJWTSessionToken + "invalid"

	_, err = postgresqlDB.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")
	if err != nil {
		t.Errorf("Failed to upload document: %v", err)
	}

	// Send a POST request to the server with the JWT token
	document := map[string]string{
		"document_name": "Test Document",
	}
	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/documents/deletedocument", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	documentFromDb, err := postgresqlDB.GetDocument(cfg.TestUserID, "Test Document")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}

	documentFromDbContent := documentFromDb.Content
	if documentFromDbContent != "Test Content" {
		t.Errorf("Expected document content %s, got %s", "Test Content", documentFromDbContent)
	}

	// Cleanup the database
	t.Cleanup(func() {
		testconfig := config.NewServerConfig()
		UserID := testconfig.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
		postgresqlDB.DeleteDocument(UserID, "Test Document")
	})
}

func TestDocumentsUpdateDocumentNameHandlerIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	err := createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	jwtToken := cfg.TestJWTSessionToken

	documentFromUpload, err := documentService.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")

	document := map[string]string{
		"documentID":        documentFromUpload.DocumentUUID.String(),
		"new_document_name": documentFromUpload.DocumentName + " Updated",
	}

	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/documents/update_document_name", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	document_response, err := documentService.GetDocument(cfg.TestUserID, "Test Document Updated")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}
	if document_response.DocumentName != "Test Document Updated" {
		t.Errorf("Expected document content %s, got %s", "Test Content Updated", document_response.Content)
	}
	if document_response.Content != "Test Content" {
		t.Errorf("Expected document content %s, got %s", "Test Content", document_response.Content)
	}

	documentUpdateContent := map[string]string{
		"documentID":           documentFromUpload.DocumentUUID.String(),
		"new_document_content": documentFromUpload.Content + " Updated",
	}

	body, _ = json.Marshal(documentUpdateContent)
	req, _ = http.NewRequest(http.MethodPut, server.URL+"/documents/update_document_content", bytes.NewBuffer(body))
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

	document_response_updated_content, err := documentService.GetDocument(cfg.TestUserID, "Test Document Updated")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}
	if document_response_updated_content.Content != "Test Content Updated" {
		t.Errorf("Expected document content %s, got %s", "Test Content Updated", document_response_updated_content.Content)
	}
	if document_response_updated_content.DocumentName != "Test Document Updated" {
		t.Errorf("Expected document content %s, got %s", "Test Document Updated", document_response_updated_content.DocumentName)
	}

	// Cleanup the database
	t.Cleanup(func() {
		UserID := cfg.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
	})
}

func TestDocumentsUpdateDocumentHandlerUnauthenticatedIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	err := createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	jwtToken := cfg.TestJWTSessionToken + "invalid"

	documentFromUpload, err := documentService.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")

	document := map[string]string{
		"documentID":        documentFromUpload.DocumentUUID.String(),
		"new_document_name": documentFromUpload.DocumentName + " Updated",
	}

	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/documents/update_document_name", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	_, err = documentService.GetDocument(cfg.TestUserID, "Test Document")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}

	documentUpdateContent := map[string]string{
		"documentID":           documentFromUpload.DocumentUUID.String(),
		"new_document_content": documentFromUpload.Content + " Updated",
	}

	body, _ = json.Marshal(documentUpdateContent)
	req, _ = http.NewRequest(http.MethodPut, server.URL+"/documents/update_document_content", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	document_response_updated_content, err := documentService.GetDocument(cfg.TestUserID, "Test Document Updated")
	if err == nil {
		t.Errorf("Should have failed to get document: %v", err)
	}
	document_response_updated_content, err = documentService.GetDocument(cfg.TestUserID, "Test Document")
	if document_response_updated_content.Content == "Test Content Updated" {
		t.Errorf("Expected document content %s, got %s", "Test Content", document_response_updated_content.Content)
	}

	// Cleanup the database
	t.Cleanup(func() {
		UserID := cfg.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
	})
}

func TestDocumentsUpdateDocumentNotMyDocumentHandlerIntegration(t *testing.T) {
	setup := SetupTestEnvironment(t)
	cfg := setup.Config
	documentService := setup.DocumentService
	postgresqlDB := setup.PostgresqlDB
	clerkInstance := setup.ClerkInstance

	err := createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentService, clerkInstance)
	server := httptest.NewServer(mux)
	defer server.Close()

	jwtToken := cfg.TestJWTSessionToken

	documentFromUpload, err := documentService.UploadDocument(cfg.TestUserID, "Test Document", "Test Content")

	document := map[string]string{
		"documentID":        documentFromUpload.DocumentUUID.String(),
		"new_document_name": documentFromUpload.DocumentName + " Updated",
	}

	body, _ := json.Marshal(document)
	req, _ := http.NewRequest(http.MethodPut, server.URL+"/documents/update_document_name", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	document_response, err := documentService.GetDocument(cfg.TestUserID, "Test Document Updated")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}
	if document_response.DocumentName != "Test Document Updated" {
		t.Errorf("Expected document content %s, got %s", "Test Content Updated", document_response.Content)
	}
	if document_response.Content != "Test Content" {
		t.Errorf("Expected document content %s, got %s", "Test Content", document_response.Content)
	}

	documentUpdateContent := map[string]string{
		"documentID":           documentFromUpload.DocumentUUID.String(),
		"new_document_content": documentFromUpload.Content + " Updated",
	}

	body, _ = json.Marshal(documentUpdateContent)
	req, _ = http.NewRequest(http.MethodPut, server.URL+"/documents/update_document_content", bytes.NewBuffer(body))
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

	document_response_updated_content, err := documentService.GetDocument(cfg.TestUserID, "Test Document Updated")
	if err != nil {
		t.Errorf("Failed to get document: %v", err)
	}
	if document_response_updated_content.Content != "Test Content Updated" {
		t.Errorf("Expected document content %s, got %s", "Test Content Updated", document_response_updated_content.Content)
	}
	if document_response_updated_content.DocumentName != "Test Document Updated" {
		t.Errorf("Expected document content %s, got %s", "Test Document Updated", document_response_updated_content.DocumentName)
	}

	// Cleanup the database
	t.Cleanup(func() {
		UserID := cfg.TestUserID
		postgresqlDB.DeleteUserInUsersTable(UserID)
	})
}
