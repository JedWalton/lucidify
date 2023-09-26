// // go:build integration
// // +build integration
package clerkapi

import (
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"testing"
)

func TestIntegration_clerk_handlers(t *testing.T) {
	testconfig := config.NewServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey
	testEmail := "clerk_handler_integration@example.com"
	firstName := "clerk_handler_firstname"
	lastName := "clerk_handler_lastname"
	password := "$sswordoatnsu28348ckj"

	storeInstance, err := store.NewStore(testconfig.PostgresqlURL)
	if err != nil {
		t.Errorf("Failed to create test store: %v", err)
	}

	userID, err := store.CreateUserInClerk(clerkSecretKey, firstName, lastName, testEmail, password)
	if err != nil {
		t.Errorf("User not created in Clerk. Reason: %v", err)
	}

	t.Cleanup(func() {
		log.Printf("Cleaning up test user: %v", userID)
		err = store.DeleteUserInClerk(clerkSecretKey, userID)
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
	err = store.UpdateUserInClerk(clerkSecretKey, userID, newFirstName, newLastName)
	if err != nil {
		t.Errorf("Failed to update user in Clerk: %v", err)
	}

	err = storeInstance.CheckUserHasExpectedFirstNameAndLastNameInUsersTable(userID, 10, newFirstName, newLastName)
	if err != nil {
		t.Errorf("User first name and last name not updated in users table: %v", err)
	}
}
