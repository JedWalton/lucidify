// //go:build integration
// // +build integration
package store

import (
	"log"
	"lucidify-api/modules/config"
	"testing"
)

func TestIntegration_usercreatedevent(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	clerkSecretKey := testconfig.ClerkSecretKey

	exists, err := CheckIfUserExists(clerkSecretKey, "clerk_handler_uce_integration@example.com")
	if err != nil {
		t.Fatalf("Failed to check if user exists in Clerk: %v", err)
	}
	if exists {
		t.Fatalf("User was not created in Clerk as it already exists.")
	}
	userID, err := CreateUserInClerk(clerkSecretKey, "clerk_handler_uce_int_firstname", "clerk_handler_uce_int_firstname", "clerk_handler_uce_integration@example.com", "$sswordoatnsu28348ckj")
	if err != nil {
		t.Fatalf("Failed to create user in Clerk: %v", err)
	}
	// Cleanup
	t.Cleanup(func() {
		if userID != "" {
			err = DeleteUserInClerk(clerkSecretKey, userID)
			if err != nil {
				log.Printf("Did not delete test user in clerk: %v\n", err)
			}
		}
	})
}

// func TestIntegration_DeleteUserInClerk(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	clerkSecretKey := testconfig.ClerkSecretKey
// 	firstName := "TestFirstName"
// 	lastName := "TestLastName"
// 	email := "testDeleteUserInClerk@example.com"
// 	password := "oeuth34c4293aoeu"
//
// 	userID, err := CreateUserInClerk(clerkSecretKey, firstName, lastName, email, password)
// 	if err != nil {
// 		t.Fatalf("Failed to create user in Clerk: %v", err)
// 	}
//
// 	storeInstance, err := NewStore(testconfig.PostgresqlURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	storeInstance.CheckUserDeletedInUsersTable(userID, 5)
//
// 	t.Cleanup(func() {
// 		err = DeleteUserInClerk(clerkSecretKey, userID)
// 		if err != nil {
// 			log.Printf("Did not delete test user: %v\n", err)
// 		}
// 		err = storeInstance.DeleteUserInUsersTable(userID)
// 		if err != nil {
// 			log.Printf("Did not delete test user in users table: %v\n", err)
// 		}
// 	})
// }
//
// func TestIntegration_UpdateUserInClerk(t *testing.T) {
// 	testconfig := config.NewTestServerConfig()
// 	clerkSecretKey := testconfig.ClerkSecretKey
// 	firstName := "TestFirstName"
// 	lastName := "TestLastName"
// 	email := "testUpdateUserInClerk@example.com"
// 	password := "soaenuth4yg8fdbioea"
//
// 	userID, err := CreateUserInClerk(clerkSecretKey, firstName, lastName, email, password)
// 	if err != nil {
// 		t.Fatalf("Failed to create user in Clerk: %v", err)
// 	}
//
// 	newFirstName := "UpdatedFirstName"
// 	newLastName := "UpdatedLastName"
// 	err = UpdateUserInClerk(clerkSecretKey, userID, newFirstName, newLastName)
// 	if err != nil {
// 		t.Fatalf("Failed to update user in Clerk: %v", err)
// 	}
//
// 	storeInstance, err := NewStore(testconfig.PostgresqlURL)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	t.Cleanup(func() {
// 		err = DeleteUserInClerk(clerkSecretKey, userID)
// 		if err != nil {
// 			log.Printf("Did not delete test user: %v\n", err)
// 		}
// 		err = storeInstance.DeleteUserInUsersTable(userID)
// 		if err != nil {
// 			log.Printf("Did not delete test user in users table: %v\n", err)
// 		}
// 	})
// }
