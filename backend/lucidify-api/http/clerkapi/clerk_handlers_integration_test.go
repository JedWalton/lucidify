// // go:build integration
// // +build integration
package clerkapi

import (
	"fmt"
	"log"
	"lucidify-api/data/store/postgresqlclient"
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

	weaviate, err := weaviateclient.NewWeaviateClientTest()
	if err != nil {
		t.Errorf("Failed to create WeaviateClient: %v", err)
	}

	postgre, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create PostgreSQLClient: %v", err)
	}
	userService, err := userservice.NewUserService(postgre, weaviate)
	if err != nil {
		t.Errorf("Failed to create UserService: %v", err)
	}

	t.Cleanup(func() {
		log.Printf("Cleaning up test user: %v", userID)
		err = clerkservice.DeleteUserInClerk(clerkSecretKey, userID)
		if err != nil {
			t.Errorf("Failed to delete test user in clerk: %v\n", err)
		}
		if userService.HasUserBeenDeleted(userID, 10) {
			t.Errorf("Failed to delete test user in users table: %v\n", err)
		}
	})

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
	// Use a channel to communicate the result of the async operation.
	resultChan := make(chan bool)
	go func() {
		// This is now running in a separate goroutine.
		for i := 0; i < 10; i++ {
			user, err := userService.GetUser(userID)
			if err == nil && user.FirstName == updatedFirstName && user.LastName == updatedLastName {
				resultChan <- true
				return
			}
			time.Sleep(1 * time.Second)
		}
		// If the loop completes without returning, the update was not successful.
		resultChan <- false
	}()

	// Use a select statement to wait for the async operation to complete or timeout.
	select {
	case success := <-resultChan:
		if !success {
			t.Errorf("User first name and last name not updated in users table")
		}
	case <-time.After(11 * time.Second): // Timeout a bit longer than the sleep * retries in the goroutine.
		t.Errorf("Update check timed out")
	}
}
