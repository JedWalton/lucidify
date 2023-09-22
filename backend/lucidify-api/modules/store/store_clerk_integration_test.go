// //go:build integration
// // +build integration
package store

import (
	"log"
	"lucidify-api/modules/config"
	"testing"
)

func TestIntegration_store_clerk(t *testing.T) {
	// Test configuration
	testconfig := config.NewTestServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey
	testEmail := "clerk_handler_uce_integration@example.com"
	firstName := "clerk_handler_uce_int_firstname"
	lastName := "clerk_handler_uce_int_lastname" // Assuming you meant to have a different last name
	password := "$sswordoatnsu28348ckj"

	_, err := CreateUserInClerk(clerkSecretKey, firstName, lastName, testEmail, password)
	if err != nil {
		log.Printf("Failed to create user in Clerk, it likely already exists so nothing to worry about: %v", err)
	}

	userID, err := getUserIDByEmail(testEmail, clerkSecretKey)
	if err != nil {
		t.Fatalf("Error getting user by email: %v", err)
	}

	newFirstName := "UpdatedFirstName"
	newLastName := "UpdatedLastName"
	err = UpdateUserInClerk(clerkSecretKey, userID, newFirstName, newLastName)
	if err != nil {
		t.Fatalf("Failed to update user in Clerk: %v", err)
	}

	// Retrieve the user to verify the update
	user, err := RetrieveUser(clerkSecretKey, userID)
	if err != nil {
		t.Fatalf("Failed to retrieve user from Clerk: %v", err)
	}

	// Check if the first name and last name match the updated values
	if user["first_name"] != newFirstName || user["last_name"] != newLastName {
		t.Fatalf("User update failed. Expected first name: %s, last name: %s. Got first name: %s, last name: %s",
			newFirstName, newLastName, user["first_name"], user["last_name"])
	}

	t.Cleanup(func() {
		log.Printf("Cleaning up test user: %v", userID)
		err = DeleteUserInClerk(clerkSecretKey, userID)
		if err != nil {
			t.Fatalf("Failed to delete test user in clerk: %v\n", err)
		}
	})
}
