// //go:build integration
// // +build integration
package syncservice

import (
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/service/userservice"
	"testing"
)

func createTestUserInDb() error {
	db, err := postgresqlclient.NewPostgreSQL()

	// the user id registered by the jwt token must exist in the local database
	user := storemodels.User{
		UserID:           "TestUserIDSyncService",
		ExternalID:       "TestCreateUserInUsersTableExternalIDSyncService",
		Username:         "TestCreateUserInUsersTableUsernameSyncService",
		PasswordEnabled:  true,
		Email:            "TestCreateUserInUsersTableSyncService@example.com",
		FirstName:        "TestCreateUserInUsersTableSyncService",
		LastName:         "TestCreateUserInUsersTableSyncService",
		ImageURL:         "https://TestCreateUserInUsersTableSyncService.com/image.jpg",
		ProfileImageURL:  "https://TestCreateUserInUsersTableSyncService.com/profile.jpg",
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

func TestSyncServiceIntegration(t *testing.T) {
	testUserID := "TestUserIDSyncService"

	err := createTestUserInDb()
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	postgre, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		log.Fatalf("Failed to create PostgreSQLClient: %v", err)
	}

	weaviate, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		log.Fatalf("Failed to create WeaviateClient: %v", err)
	}
	userService, err := userservice.NewUserService(postgre, weaviate)
	if err != nil {
		log.Fatalf("Failed to create UserService: %v", err)
	}

	defer userService.DeleteUser(testUserID) // Cleanup the test user after the test

	// Initialize SyncService
	syncSrv, err := NewSyncService()
	if err != nil {
		t.Fatalf("Failed to initialize SyncService: %v", err)
	}

	testKeys := []string{"conversationHistory", "prompts", "folders"}
	testValue := "This is a test value."

	for _, testKey := range testKeys {
		runTestForKey(t, syncSrv, testUserID, testKey, testValue)
	}
}

func runTestForKey(t *testing.T, syncSrv SyncService, userID, key, value string) {
	// Test HandleSet
	resp := syncSrv.HandleSet(userID, key, value)
	if !resp.Success {
		t.Fatalf("HandleSet failed for key '%s': %s", key, resp.Message)
	}

	// Test HandleGet
	resp = syncSrv.HandleGet(userID, key)
	if !resp.Success {
		t.Fatalf("HandleGet failed for key '%s': %s", key, resp.Message)
	}
	if resp.Data.(string) != value {
		t.Fatalf("For key '%s', expected data to be '%s' but got '%s'", key, value, resp.Data)
	}

	// Test HandleClearConversations
	resp = syncSrv.HandleClearConversations(userID)
	if !resp.Success {
		t.Fatalf("HandleClearConversations failed for key '%s': %s", key, resp.Message)
	}

	// Verify conversationHistory & folders are cleared
	if key == "conversationHistory" || key == "folders" {
		// Verify data is cleared
		resp = syncSrv.HandleGet(userID, key)
		if resp.Success {
			t.Fatalf("For key '%s', expected data to be cleared, but HandleGet was successful", key)
		}
		return
	}
}
