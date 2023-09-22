// // go:build integration
// // +build integration
package clerk

import (
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"testing"
)

func TestIntegration_clerk_handlers(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey
	testEmail := "clerk_handler_integration@example.com"
	firstName := "clerk_handler_firstname"
	lastName := "clerk_handler_lastname"
	password := "$sswordoatnsu28348ckj"

	userID, err := store.CreateUserInClerk(clerkSecretKey, firstName, lastName, testEmail, password)
	if err != nil {
		log.Printf("Failed to create user in Clerk, it likely already exists so nothing to worry about: %v", err)
	}

	storeInstance, err := store.NewStore(testconfig.PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	err = storeInstance.CheckIfUserInUsersTable(userID, 20)
	if err != nil {
		t.Fatalf("User not found after creation: %v", err)
	}

	newFirstName := "updated_clerk_handler_firstname"
	newLastName := "updated_clerk_handler_lastname"
	err = store.UpdateUserInClerk(clerkSecretKey, userID, newFirstName, newLastName)
	if err != nil {
		t.Fatalf("Failed to update user in Clerk: %v", err)
	}

	err = storeInstance.CheckUserHasExpectedFirstNameAndLastNameInUsersTable(userID, 20, newFirstName, newLastName)
	if err != nil {
		t.Fatalf("User first name and last name not updated in users table: %v", err)
	}

	t.Cleanup(func() {
		log.Printf("Cleaning up test user: %v", userID)
		err = store.DeleteUserInClerk(clerkSecretKey, userID)
		if err != nil {
			t.Fatalf("Failed to delete test user in clerk: %v\n", err)
		}
		err = storeInstance.DeleteUserInUsersTable(userID)
		if err != nil {
			t.Fatalf("Failed to delete test user in users table: %v\n", err)
		}
	})
}
