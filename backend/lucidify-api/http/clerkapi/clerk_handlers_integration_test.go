// // go:build integration
// // +build integration
package clerkapi

import (
	"fmt"
	"log"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/server/config"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/userservice"
	"testing"
	"time"
)

func TestIntegration_clerk_handlers(t *testing.T) {
	testconfig := config.NewServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey
	testEmail := fmt.Sprintf("clerk_handler_integration_%d@example.com", time.Now().UnixNano())
	firstName := "clerk_handler_firstname"
	lastName := "clerk_handler_lastname"
	password := "$sswordoatnsu28348ckj"

	userID, err := clerkservice.CreateUserInClerk(clerkSecretKey, firstName, lastName, testEmail, password)
	if err != nil {
		t.Errorf("User not created in Clerk. Reason: %v", err)
	}

	t.Cleanup(func() {
		log.Printf("Cleaning up test user: %v", userID)
		err = clerkservice.DeleteUserInClerk(clerkSecretKey, userID)
		if err != nil {
			t.Errorf("Failed to delete test user in clerk: %v\n", err)
		}
		weaviate, err := weaviateclient.NewWeaviateClientTest()
		if err != nil {
			t.Errorf("Failed to create WeaviateClient: %v", err)
		}
		userService, err := userservice.NewUserService(weaviate)
		if err != nil {
			t.Errorf("Failed to create UserService: %v", err)
		}
		if userService.HasUserBeenDeleted(userID, 10) {
			t.Errorf("Failed to delete test user in users table: %v\n", err)
		}
	})

	weaviate, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		t.Errorf("Failed to create WeaviateClient: %v", err)
	}
	userService, err := userservice.NewUserService(weaviate)
	if err != nil {
		t.Errorf("Failed to create UserService: %v", err)
	}
	_, err = userService.GetUserWithRetries(userID, 10)
	if err != nil {
		t.Errorf("User not found after creation: %v", err)
	}

	updatedFirstName := "updated_clerk_handler_firstname"
	updatedLastName := "updated_clerk_handler_lastname"
	err = clerkservice.UpdateUserInClerk(clerkSecretKey, userID, updatedFirstName, updatedLastName)
	if err != nil {
		t.Errorf("Failed to update user in Clerk: %v", err)
	}

	var updated bool
	for i := 0; i < 10; i++ {
		user, err := userService.GetUser(userID)
		if err == nil && user.FirstName == updatedFirstName && user.LastName == updatedLastName {
			updated = true
			break
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	if !updated {
		t.Errorf("User first name and last name not updated in users table")
	}
}
