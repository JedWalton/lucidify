// //go:build integration
// // +build integration
package store

import (
	"lucidify-api/modules/config"
	"testing"
)

func TestCreateUserInUsersTable(t *testing.T) {
	testconfig := config.NewServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewStore(PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	user := User{
		UserID:           "TestCreateUserInUsersTableUserID",
		ExternalID:       "TestCreateUserInUsersTableExternalID",
		Username:         "TestCreateUserInUsersTableUsername",
		PasswordEnabled:  true,
		Email:            "TestCreateUserInUsersTable@example.com",
		FirstName:        "TestCreateUserInUsersTableCreateTest",
		LastName:         "TestCreateUserInUsersTableUser",
		ImageURL:         "https://TestCreateUserInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://TestCreateUserInUsersTable.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err = store.CreateUserInUsersTable(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Check if the user exists
	err = store.CheckIfUserInUsersTable(user.UserID, 3)
	if err != nil {
		t.Fatalf("User not found after creation: %v", err)
	}

	// Register cleanup function
	t.Cleanup(func() {
		store.DeleteUserInUsersTable(user.UserID)
	})
}

func TestUpdateUserInUsersTable(t *testing.T) {
	testconfig := config.NewServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewStore(PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Create a user first
	user := User{
		UserID:           "TestUpdateUserInUsersTableUserID",
		ExternalID:       "TestUpdateUserInUsersTableUserIDExternalID",
		Username:         "TestUpdateUserInUsersTableUserIDUsername",
		PasswordEnabled:  true,
		Email:            "TestUpdateUserInUsersTableUserID@example.com",
		FirstName:        "TestUpdateUserInUsersTableUserID",
		LastName:         "User",
		ImageURL:         "https://updateTest.com/image.jpg",
		ProfileImageURL:  "https://updateTest.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err = store.CreateUserInUsersTable(user)
	if err != nil {
		t.Fatalf("Failed to create user for update test: %v", err)
	}

	// Update the user
	user.FirstName = "UpdatedFirstName"
	user.LastName = "UpdatedLastName"
	err = store.UpdateUserInUsersTable(user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Check if the user has the expected first name and last name
	err = store.CheckUserHasExpectedFirstNameAndLastNameInUsersTable(user.UserID, 3, "UpdatedFirstName", "UpdatedLastName")
	if err != nil {
		t.Fatalf("User not updated correctly: %v", err)
	}

	// Register cleanup function
	t.Cleanup(func() {
		store.DeleteUserInUsersTable(user.UserID)
	})

}

func TestGetUserInUsersTable(t *testing.T) {
	testconfig := config.NewServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewStore(PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Create a user first
	user := User{
		UserID:           "TestGetUserInUsersTableUserID",
		ExternalID:       "TestGetUserInUsersTableExternalID",
		Username:         "TestGetUserInUsersTableUsername",
		PasswordEnabled:  true,
		Email:            "TestGetUserInUsersTable@example.com",
		FirstName:        "GetTest",
		LastName:         "User",
		ImageURL:         "https://TestGetUserInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://getTest.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err = store.CreateUserInUsersTable(user)
	if err != nil {
		t.Fatalf("Failed to create user for get test: %v", err)
	}

	// Fetch the user
	fetchedUser, err := store.GetUserInUsersTable(user.UserID)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	if fetchedUser.UserID != user.UserID || fetchedUser.Email != user.Email {
		t.Fatalf("Fetched user does not match created user. Expected user ID: %v, got: %v. Expected email: %v, got: %v", user.UserID, fetchedUser.UserID, user.Email, fetchedUser.Email)
	}

	// Cleanup
	t.Cleanup(func() {
		store.DeleteUserInUsersTable(user.UserID)
	})
}

func TestDeleteUserInUsersTable(t *testing.T) {
	testconfig := config.NewServerConfig()
	PostgresqlURL := testconfig.PostgresqlURL

	store, err := NewStore(PostgresqlURL)
	if err != nil {
		t.Fatalf("Failed to create test store: %v", err)
	}

	// Create a user first
	user := User{
		UserID:           "TestDeleteUserInUsersTableUserID",
		ExternalID:       "TestDeleteUserInUsersTableExternalID",
		Username:         "TestDeleteUserInUsersTableUsername",
		PasswordEnabled:  true,
		Email:            "TestDeleteUserInUsersTable@example.com",
		FirstName:        "TestDeleteUserInUsersTable",
		LastName:         "TestDeleteUserInUsersTableUser",
		ImageURL:         "https://TestDeleteUserInUsersTable.com/image.jpg",
		ProfileImageURL:  "https://TestDeleteUserInUsersTable.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err = store.CreateUserInUsersTable(user)
	if err != nil {
		t.Fatalf("Failed to create user for delete test: %v", err)
	}

	// Delete the user
	err = store.DeleteUserInUsersTable(user.UserID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Check if the user has been deleted
	err = store.CheckUserDeletedInUsersTable(user.UserID, 3)
	if err != nil {
		t.Fatalf("User still exists after deletion: %v", err)
	}

	t.Cleanup(func() {
		store.DeleteUserInUsersTable(user.UserID)
	})
}
