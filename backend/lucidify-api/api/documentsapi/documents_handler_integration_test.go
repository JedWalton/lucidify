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
	postgresqlclient2 "lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/store"
	"lucidify-api/modules/store/storemodels"
	"lucidify-api/modules/store/weaviateclient"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createTestUserInDb() error {
	testconfig := config.NewServerConfig()
	db, err := postgresqlclient2.NewPostgreSQL()

	// the user id registered by the jwt token must exist in the local database
	user := postgresqlclient2.User{
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
	db, err := postgresqlclient2.NewPostgreSQL()

	// the user id registered by the jwt token must exist in the local database
	user := postgresqlclient2.User{
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

func TestDocumentsUploadHandlerIntegration(t *testing.T) {
	cfg := config.NewServerConfig()
	postgresqlDB, err := postgresqlclient2.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	// Setup the real environment
	clerkInstance, err := clerkclient.NewClerkClient(cfg.ClerkSecretKey)
	if err != nil {
		t.Errorf("Failed to create Clerk client: %v", err)
	}
	weaviateDB, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		t.Errorf("Failed to create Weaviate client: %v", err)
	}
	err = createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}
	documentsService := store.NewDocumentService(postgresqlDB, weaviateDB)

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentsService, clerkInstance)
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
	cfg := config.NewServerConfig()
	postgresqlDB, err := postgresqlclient2.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	// Setup the real environment
	clerkInstance, err := clerkclient.NewClerkClient(cfg.ClerkSecretKey)
	if err != nil {
		t.Errorf("Failed to create Clerk client: %v", err)
	}
	weaviateDB, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		t.Errorf("Failed to create Weaviate client: %v", err)
	}
	err = createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}
	documentsService := store.NewDocumentService(postgresqlDB, weaviateDB)

	// Create a test server
	mux := http.NewServeMux()
	SetupRoutes(cfg, mux, documentsService, clerkInstance)
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
	cfg := config.NewServerConfig()
	postgresqlDB, err := postgresqlclient2.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	// Setup the real environment
	clerkInstance, err := clerkclient.NewClerkClient(cfg.ClerkSecretKey)
	if err != nil {
		t.Errorf("Failed to create Clerk client: %v", err)
	}
	weaviateDB, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		t.Errorf("Failed to create Weaviate client: %v", err)
	}
	err = createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}
	documentService := store.NewDocumentService(postgresqlDB, weaviateDB)

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
	cfg := config.NewServerConfig()
	postgresqlDB, err := postgresqlclient2.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	// Setup the real environment
	clerkInstance, err := clerkclient.NewClerkClient(cfg.ClerkSecretKey)
	if err != nil {
		t.Errorf("Failed to create Clerk client: %v", err)
	}
	weaviateDB, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		t.Errorf("Failed to create Weaviate client: %v", err)
	}
	err = createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}
	documentService := store.NewDocumentService(postgresqlDB, weaviateDB)

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

//	func TestDocumentsGetAllDocumentsHandlerIntegration(t *testing.T) {
//		testconfig := config.NewServerConfig()
//		db, err := postgresqlclient2.NewPostgreSQL()
//
//		clerkInstance, err := clerkclient.NewClerkClient(testconfig.ClerkSecretKey)
//
//		createTestUserInDb()
//
//		if err != nil {
//			t.Errorf("Failed to create Clerk client: %v", err)
//		}
//		cfg := &config.ServerConfig{}
//
//		// Create a test server
//		mux := http.NewServeMux()
//		SetupRoutes(cfg, mux, db, clerkInstance)
//		server := httptest.NewServer(mux)
//		defer server.Close()
//
//		jwtToken := testconfig.TestJWTSessionToken
//
//		db.UploadDocument(testconfig.TestUserID, "Test Document", "Test Content")
//		db.UploadDocument(testconfig.TestUserID, "Test Document 2", "Test Content 2")
//		db.UploadDocument(testconfig.TestUserID, "Test Document 3", "Test Content 3")
//
//		req, _ := http.NewRequest(http.MethodGet, server.URL+"/documents/getalldocuments", nil)
//		req.Header.Set("Authorization", "Bearer "+jwtToken)
//		client := &http.Client{}
//		resp, err := client.Do(req)
//		if err != nil {
//			t.Errorf("Failed to send request: %v", err)
//		}
//		defer resp.Body.Close()
//
//		// Check the response
//		if resp.StatusCode != http.StatusOK {
//			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
//		}
//
//		// Read the response body
//		respBody, err := io.ReadAll(resp.Body)
//		if err != nil {
//			t.Errorf("Failed to read response body: %v", err)
//		}
//
//		// Unmarshal the response body into a slice of Document objects
//		var respDocuments []storemodels.Document
//		err = json.Unmarshal(respBody, &respDocuments)
//		if err != nil {
//			t.Errorf("Failed to unmarshal response body: %v", err)
//		}
//
//		// Check if the returned documents are correct
//		if len(respDocuments) != 3 {
//			t.Errorf("Expected 3 documents, got %d", len(respDocuments))
//		}
//
//		expectedDocs := []string{"Test Document", "Test Document 2", "Test Document 3"}
//		for i, doc := range respDocuments {
//			if doc.DocumentName != expectedDocs[i] {
//				t.Errorf("Expected document name %s, got %s", expectedDocs[i], doc.DocumentName)
//			}
//		}
//
//		// Cleanup the database
//		t.Cleanup(func() {
//			testconfig := config.NewServerConfig()
//			UserID := testconfig.TestUserID
//			db.DeleteUserInUsersTable(UserID)
//			db.DeleteDocument(UserID, "Test Document")
//			db.DeleteDocument(UserID, "Test Document 2")
//			db.DeleteDocument(UserID, "Test Document 3")
//		})
//	}
//
//	func TestDocumentsGetAllDocumentsHandlerUnauthenticatedIntegration(t *testing.T) {
//		testconfig := config.NewServerConfig()
//		db, err := postgresqlclient2.NewPostgreSQL()
//
//		clerkInstance, err := clerkclient.NewClerkClient(testconfig.ClerkSecretKey)
//
//		createTestUserInDb()
//
//		if err != nil {
//			t.Errorf("Failed to create Clerk client: %v", err)
//		}
//		cfg := &config.ServerConfig{}
//
//		// Create a test server
//		mux := http.NewServeMux()
//		SetupRoutes(cfg, mux, db, clerkInstance)
//		server := httptest.NewServer(mux)
//		defer server.Close()
//
//		jwtToken := testconfig.TestJWTSessionToken + "invalid"
//
//		db.UploadDocument(testconfig.TestUserID, "Test Document", "Test Content")
//		db.UploadDocument(testconfig.TestUserID, "Test Document 2", "Test Content 2")
//		db.UploadDocument(testconfig.TestUserID, "Test Document 3", "Test Content 3")
//
//		req, _ := http.NewRequest(http.MethodGet, server.URL+"/documents/getalldocuments", nil)
//		req.Header.Set("Authorization", "Bearer "+jwtToken)
//		client := &http.Client{}
//		resp, err := client.Do(req)
//		if err != nil {
//			t.Errorf("Failed to send request: %v", err)
//		}
//		defer resp.Body.Close()
//
//		// Check the response
//		if resp.StatusCode != http.StatusBadRequest {
//			t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
//		}
//
//		// Cleanup the database
//		t.Cleanup(func() {
//			testconfig := config.NewServerConfig()
//			UserID := testconfig.TestUserID
//			db.DeleteUserInUsersTable(UserID)
//			db.DeleteDocument(UserID, "Test Document")
//			db.DeleteDocument(UserID, "Test Document 2")
//			db.DeleteDocument(UserID, "Test Document 3")
//		})
//	}
//
//	func TestDocumentsGetAllDocumentsHandlerUnauthenticatedOtherUserIntegration(t *testing.T) {
//		testconfig := config.NewServerConfig()
//		db, err := postgresqlclient2.NewPostgreSQL()
//
//		clerkInstance, err := clerkclient.NewClerkClient(testconfig.ClerkSecretKey)
//
//		createTestUserInDb()
//
//		UserID2 := createASecondTestUserInDb()
//
//		if err != nil {
//			t.Errorf("Failed to create Clerk client: %v", err)
//		}
//		cfg := &config.ServerConfig{}
//
//		// Create a test server
//		mux := http.NewServeMux()
//		SetupRoutes(cfg, mux, db, clerkInstance)
//		server := httptest.NewServer(mux)
//		defer server.Close()
//
//		jwtToken := testconfig.TestJWTSessionToken
//
//		db.UploadDocument(UserID2, "Test Document", "Test Content")
//		db.UploadDocument(UserID2, "Test Document 2", "Test Content 2")
//		db.UploadDocument(UserID2, "Test Document 3", "Test Content 3")
//
//		req, _ := http.NewRequest(http.MethodGet, server.URL+"/documents/getalldocuments", nil)
//		req.Header.Set("Authorization", "Bearer "+jwtToken)
//		client := &http.Client{}
//		resp, err := client.Do(req)
//		if err != nil {
//			t.Errorf("Failed to send request: %v", err)
//		}
//		defer resp.Body.Close()
//
//		// Check the response
//		if resp.StatusCode != http.StatusOK {
//			t.Errorf("Expected status code %v, got %v", http.StatusOK, resp.StatusCode)
//		}
//
//		// Read the response body
//		respBody, err := io.ReadAll(resp.Body)
//		if err != nil {
//			t.Errorf("Failed to read response body: %v", err)
//		}
//
//		// Unmarshal the response body into a slice of Document objects
//		var respDocuments []storemodels.Document
//		err = json.Unmarshal(respBody, &respDocuments)
//		if err != nil {
//			t.Errorf("Failed to unmarshal response body: %v", err)
//		}
//
//		// Check if the returned documents are correct
//		if len(respDocuments) == 3 {
//			t.Errorf("Expected 0 documents, got %d", len(respDocuments))
//		}
//
//		expectedDocs := []string{"Test Document", "Test Document 2", "Test Document 3"}
//		for i, doc := range respDocuments {
//			if doc.DocumentName == expectedDocs[i] {
//				t.Errorf("Expected to not mach document name %s, got %s", expectedDocs[i], doc.DocumentName)
//			}
//		}
//		// Cleanup the database
//		t.Cleanup(func() {
//			db.DeleteUserInUsersTable(UserID2)
//			db.DeleteDocument(UserID2, "Test Document")
//			db.DeleteDocument(UserID2, "Test Document 2")
//			db.DeleteDocument(UserID2, "Test Document 3")
//		})
//	}
func TestDocumentsDeleteDocumentHandlerIntegration(t *testing.T) {

	cfg := config.NewServerConfig()
	postgresqlDB, err := postgresqlclient2.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	// Setup the real environment
	clerkInstance, err := clerkclient.NewClerkClient(cfg.ClerkSecretKey)
	if err != nil {
		t.Errorf("Failed to create Clerk client: %v", err)
	}
	weaviateDB, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		t.Errorf("Failed to create Weaviate client: %v", err)
	}
	err = createTestUserInDb()
	if err != nil {
		t.Errorf("Failed to create test user in db: %v", err)
	}
	documentService := store.NewDocumentService(postgresqlDB, weaviateDB)

	createTestUserInDb()

	if err != nil {
		t.Errorf("Failed to create Clerk client: %v", err)
	}

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

//
// func TestDocumentsDeleteDocumentHandlerUnauthenticatedIntegration(t *testing.T) {
// 	testconfig := config.NewServerConfig()
// 	db, err := postgresqlclient2.NewPostgreSQL()
// 	clerkInstance, err := clerkclient.NewClerkClient(testconfig.ClerkSecretKey)
// 	createTestUserInDb()
//
// 	if err != nil {
// 		t.Errorf("Failed to create Clerk client: %v", err)
// 	}
// 	cfg := &config.ServerConfig{}
//
// 	// Create a test server
// 	mux := http.NewServeMux()
// 	SetupRoutes(cfg, mux, db, clerkInstance)
// 	server := httptest.NewServer(mux)
// 	defer server.Close()
//
// 	// Obtain a JWT token from Clerk
// 	jwtToken := testconfig.TestJWTSessionToken + "invalid"
//
// 	_, err = db.UploadDocument(testconfig.TestUserID, "Test Document", "Test Content")
// 	if err != nil {
// 		t.Errorf("Failed to upload document: %v", err)
// 	}
//
// 	// Send a POST request to the server with the JWT token
// 	document := map[string]string{
// 		"document_name": "Test Document",
// 	}
// 	body, _ := json.Marshal(document)
// 	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/documents/deletedocument", bytes.NewBuffer(body))
// 	req.Header.Set("Authorization", "Bearer "+jwtToken)
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Errorf("Failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()
//
// 	// Check the response
// 	if resp.StatusCode != http.StatusBadRequest {
// 		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
// 	}
//
// 	documentFromDb, err := db.GetDocument(testconfig.TestUserID, "Test Document")
// 	if err != nil {
// 		t.Errorf("Failed to get document: %v", err)
// 	}
//
// 	documentFromDbContent := documentFromDb.Content
// 	if documentFromDbContent != "Test Content" {
// 		t.Errorf("Expected document content %s, got %s", "Test Content", documentFromDbContent)
// 	}
//
// 	// Cleanup the database
// 	t.Cleanup(func() {
// 		testconfig := config.NewServerConfig()
// 		UserID := testconfig.TestUserID
// 		db.DeleteUserInUsersTable(UserID)
// 		db.DeleteDocument(UserID, "Test Document")
// 	})
// }
//
// func TestDocumentsUpdateDocumentHandlerIntegration(t *testing.T) {
// 	testconfig := config.NewServerConfig()
// 	db, err := postgresqlclient2.NewPostgreSQL()
//
// 	clerkInstance, err := clerkclient.NewClerkClient(testconfig.ClerkSecretKey)
// 	if err != nil {
// 		t.Errorf("Failed to create Clerk client: %v", err)
// 	}
//
// 	createTestUserInDb()
// 	cfg := &config.ServerConfig{}
//
// 	// Create a test server
// 	mux := http.NewServeMux()
// 	SetupRoutes(cfg, mux, db, clerkInstance)
// 	server := httptest.NewServer(mux)
// 	defer server.Close()
//
// 	jwtToken := testconfig.TestJWTSessionToken
//
// 	document := map[string]string{
// 		"document_name": "Test Document",
// 		"content":       "Test Content Updated",
// 	}
//
// 	db.UploadDocument(testconfig.TestUserID, "Test Document", "Test Content")
//
// 	body, _ := json.Marshal(document)
// 	req, _ := http.NewRequest(http.MethodPut, server.URL+"/documents/updatedocument", bytes.NewBuffer(body))
// 	req.Header.Set("Authorization", "Bearer "+jwtToken)
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Errorf("Failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()
//
// 	// Check the response
// 	if resp.StatusCode != http.StatusOK {
// 		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
// 	}
//
// 	document_response, err := db.GetDocument(testconfig.TestUserID, "Test Document")
// 	if err != nil {
// 		t.Errorf("Failed to get document: %v", err)
// 	}
// 	if document_response.Content != "Test Content Updated" {
// 		t.Errorf("Expected document content %s, got %s", "Test Content Updated", document_response.Content)
// 	}
//
// 	// Cleanup the database
// 	t.Cleanup(func() {
// 		testconfig := config.NewServerConfig()
// 		UserID := testconfig.TestUserID
// 		db.DeleteUserInUsersTable(UserID)
// 		db.DeleteDocument(UserID, "Test Document")
// 	})
// }
//
// func TestDocumentsUpdateDocumentHandlerUnauthorizedIntegration(t *testing.T) {
// 	testconfig := config.NewServerConfig()
// 	db, err := postgresqlclient2.NewPostgreSQL()
//
// 	clerkInstance, err := clerkclient.NewClerkClient(testconfig.ClerkSecretKey)
// 	if err != nil {
// 		t.Errorf("Failed to create Clerk client: %v", err)
// 	}
//
// 	createTestUserInDb()
// 	cfg := &config.ServerConfig{}
//
// 	// Create a test server
// 	mux := http.NewServeMux()
// 	SetupRoutes(cfg, mux, db, clerkInstance)
// 	server := httptest.NewServer(mux)
// 	defer server.Close()
//
// 	jwtToken := testconfig.TestJWTSessionToken + "invalid"
//
// 	document := map[string]string{
// 		"document_name": "Test Document",
// 		"content":       "Test Content Updated",
// 	}
//
// 	db.UploadDocument(testconfig.TestUserID, "Test Document", "Test Content")
//
// 	body, _ := json.Marshal(document)
// 	req, _ := http.NewRequest(http.MethodPut, server.URL+"/documents/updatedocument", bytes.NewBuffer(body))
// 	req.Header.Set("Authorization", "Bearer "+jwtToken)
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Errorf("Failed to send request: %v", err)
// 	}
// 	defer resp.Body.Close()
//
// 	// Check the response
// 	if resp.StatusCode != http.StatusBadRequest {
// 		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
// 	}
//
// 	document_response, err := db.GetDocument(testconfig.TestUserID, "Test Document")
// 	if err != nil {
// 		t.Errorf("Failed to get document: %v", err)
// 	}
// 	if document_response.Content == "Test Content Updated" {
// 		t.Errorf("Expected document content to have not been updated.")
// 	}
//
// 	// Cleanup the database
// 	t.Cleanup(func() {
// 		testconfig := config.NewServerConfig()
// 		UserID := testconfig.TestUserID
// 		db.DeleteUserInUsersTable(UserID)
// 		db.DeleteDocument(UserID, "Test Document")
// 	})
// }
