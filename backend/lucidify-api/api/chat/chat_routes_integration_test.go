//go:build integration
// +build integration

package chat

import (
	"io/ioutil"
	"log"
	"lucidify-api/modules/clerkclient"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/store"
	"lucidify-api/modules/store/weaviateclient"
	"os/exec"
	"testing"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func WriteToFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func ReadFromFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func TestIntegration_unauthorized_chat(t *testing.T) {
	MakeCurlRequest := func() (string, error) {
		cmd := exec.Command("curl", "-s", "-X", "POST", "http://localhost:8080/chat", "-H", "Content-Type: application/json", "-d", `{"message": "hello"}`)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	// fileContent, err := ReadFromFile("chat.txt")
	// if err != nil {
	// 	t.Errorf("Failed to read from file: %v", err)
	// }

	response, err := MakeCurlRequest()
	if err != nil {
		t.Errorf("Failed to make curl request: %v", err)
	}
	expectedResponse := "couldn't find cookie __session"

	if response != expectedResponse {
		t.Errorf("Expected %q but got %q", expectedResponse, response)
	}
}

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

func TestIntegration_chat(t *testing.T) {
	// 1. Set up the test environment and get the JWT token
	setup := SetupTestEnvironment(t)
	jwtToken := setup.Config.TestJWTSessionToken

	MakeCurlRequest := func() (string, error) {
		// 2. Include the JWT token in the request headers
		cmd := exec.Command("curl", "-s", "-X", "POST", "http://localhost:8080/chat",
			"-H", "Content-Type: application/json",
			"-H", "Authorization: Bearer "+jwtToken,
			"-d", `{"message": "hello"}`)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	fileContent, err := ReadFromFile("chat.txt")
	if err != nil {
		t.Errorf("Failed to read from file: %v", err)
	}

	response, err := MakeCurlRequest()
	if err != nil {
		t.Errorf("Failed to make curl request: %v", err)
	}

	if response != fileContent {
		t.Errorf("Expected %q but got %q", fileContent, response)
	}
}
