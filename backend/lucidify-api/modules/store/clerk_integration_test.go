// //go:build integration
// // +build integration
package store

import (
	"lucidify-api/modules/config"
	"testing"
)

func TestIntegration_CreateUserInClerk(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey
	firstName := "TestFirstName"
	lastName := "TestLastName"
	email := "testCreateUserInClerk@example.com"
	password := "823f458ide7012oeuC.p,p"

	userID, err := CreateUserInClerk(clerkSecretKey, firstName, lastName, email, password)
	if err != nil {
		t.Fatalf("Failed to create user in Clerk: %v", err)
	}

	// Cleanup
	err = DeleteUserInClerk(clerkSecretKey, userID)
	if err != nil {
		t.Fatalf("Failed to delete test user: %v", err)
	}
}

func TestIntegration_DeleteUserInClerk(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey
	firstName := "TestFirstName"
	lastName := "TestLastName"
	email := "testDeleteUserInClerk@example.com"
	password := "oeuth34c4293aoeu"

	userID, err := CreateUserInClerk(clerkSecretKey, firstName, lastName, email, password)
	if err != nil {
		t.Fatalf("Failed to create user in Clerk: %v", err)
	}

	err = DeleteUserInClerk(clerkSecretKey, userID)
	if err != nil {
		t.Fatalf("Failed to delete user in Clerk: %v", err)
	}
}

func TestIntegration_UpdateUserInClerk(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey
	firstName := "TestFirstName"
	lastName := "TestLastName"
	email := "testUpdateUserInClerk@example.com"
	password := "soaenuth4yg8fdbioea"

	userID, err := CreateUserInClerk(clerkSecretKey, firstName, lastName, email, password)
	if err != nil {
		t.Fatalf("Failed to create user in Clerk: %v", err)
	}

	newFirstName := "UpdatedFirstName"
	newLastName := "UpdatedLastName"
	err = UpdateUserInClerk(clerkSecretKey, userID, newFirstName, newLastName)
	if err != nil {
		t.Fatalf("Failed to update user in Clerk: %v", err)
	}

	// Cleanup
	err = DeleteUserInClerk(clerkSecretKey, userID)
	if err != nil {
		t.Fatalf("Failed to delete test user: %v", err)
	}
}
