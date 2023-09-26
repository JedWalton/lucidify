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
	testconfig := config.NewServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey
	testEmail := "store_clerk_integration@example.com"
	firstName := "store_clerk_firstname"
	lastName := "store_clerk_lastname"
	password := "$sswordoatnsu28348ckj"

	_, err := CreateUserInClerk(clerkSecretKey, firstName, lastName, testEmail, password)
	if err != nil {
		log.Printf("User not created in Clerk, likely already exists.")
	}

	userID, err := GetUserIDByEmail(testEmail, clerkSecretKey)
	if err != nil {
		t.Errorf("Error getting user by email: %v", err)
	}

	newFirstName := "updated_store_clerk_firstname"
	newLastName := "updated_store_clerk_lastname"
	err = UpdateUserInClerk(clerkSecretKey, userID, newFirstName, newLastName)
	if err != nil {
		t.Errorf("Failed to update user in Clerk: %v", err)
	}

	// Retrieve the user to verify the update
	user, err := RetrieveUser(clerkSecretKey, userID)
	if err != nil {
		t.Errorf("Failed to retrieve user from Clerk: %v", err)
	}

	// Check if the first name and last name match the updated values
	if user["first_name"] != newFirstName || user["last_name"] != newLastName {
		t.Errorf("User update failed. Expected first name: %s, last name: %s. Got first name: %s, last name: %s",
			newFirstName, newLastName, user["first_name"], user["last_name"])
	}

	t.Cleanup(func() {
		log.Printf("Cleaning up test user: %v", userID)
		err = DeleteUserInClerk(clerkSecretKey, userID)
		if err != nil {
			t.Errorf("Failed to delete test user in clerk: %v\n", err)
		}
	})
}
