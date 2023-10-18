// // go:build integration
// // +build integration
package clerkapi

import (
	"fmt"
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/server/config"
	"lucidify-api/service/clerkservice"
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

	storeInstance, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

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
		err = storeInstance.CheckUserDeletedInUsersTable(userID, 10)
		if err != nil {
			t.Errorf("Failed to delete test user in users table: %v\n", err)
		}
	})

	err = storeInstance.CheckIfUserInUsersTable(userID, 10)
	if err != nil {
		t.Errorf("User not found after creation: %v", err)
	}

	newFirstName := "updated_clerk_handler_firstname"
	newLastName := "updated_clerk_handler_lastname"
	err = clerkservice.UpdateUserInClerk(clerkSecretKey, userID, newFirstName, newLastName)
	if err != nil {
		t.Errorf("Failed to update user in Clerk: %v", err)
	}

	var userFromDb *storemodels.User
	for i := 0; i < 5; i++ {
		userFromDb, err = storeInstance.GetUserInUsersTable(userID)
		if err == nil && userFromDb.FirstName == newFirstName && userFromDb.LastName == newLastName {
			break
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	if userFromDb.FirstName != newFirstName || userFromDb.LastName != newLastName {
		t.Errorf("User first name and last name not updated in users table: %v", err)
	}
}
