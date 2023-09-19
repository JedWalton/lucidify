//go:build integration
// +build integration

package store

import (
	"lucidify-api/modules/testutils"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := testutils.SetupDB()
	defer db.Close()

	store := &Store{db: db}

	user := User{
		UserID:           "createTestUserID",
		ExternalID:       "createTestExternalID",
		Username:         "createTestUsername",
		PasswordEnabled:  true,
		Email:            "createTest@example.com",
		FirstName:        "CreateTest",
		LastName:         "User",
		ImageURL:         "https://createTest.com/image.jpg",
		ProfileImageURL:  "https://createTest.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Cleanup
	err = store.DeleteUser(user.UserID)
	if err != nil {
		t.Fatalf("Failed to clean up test user: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	db := testutils.SetupDB()
	defer db.Close()

	store := &Store{db: db}

	// Create a user first
	user := User{
		UserID:           "updateTestUserID",
		ExternalID:       "updateTestExternalID",
		Username:         "updateTestUsername",
		PasswordEnabled:  true,
		Email:            "updateTest@example.com",
		FirstName:        "UpdateTest",
		LastName:         "User",
		ImageURL:         "https://updateTest.com/image.jpg",
		ProfileImageURL:  "https://updateTest.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user for update test: %v", err)
	}

	// Update the user
	user.FirstName = "UpdatedFirstName"
	user.LastName = "UpdatedLastName"
	err = store.UpdateUser(user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}

	// Fetch the user and check if the updates are reflected
	updatedUser, err := store.GetUser(user.UserID)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	if updatedUser.FirstName != "UpdatedFirstName" || updatedUser.LastName != "UpdatedLastName" {
		t.Fatalf("User update failed. Expected first name: UpdatedFirstName, got: %v. Expected last name: UpdatedLastName, got: %v", updatedUser.FirstName, updatedUser.LastName)
	}

	// Cleanup
	err = store.DeleteUser(user.UserID)
	if err != nil {
		t.Fatalf("Failed to clean up test user: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	db := testutils.SetupDB()
	defer db.Close()

	store := &Store{db: db}

	// Create a user first
	user := User{
		UserID:           "getTestUserID",
		ExternalID:       "getTestExternalID",
		Username:         "getTestUsername",
		PasswordEnabled:  true,
		Email:            "getTest@example.com",
		FirstName:        "GetTest",
		LastName:         "User",
		ImageURL:         "https://getTest.com/image.jpg",
		ProfileImageURL:  "https://getTest.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user for get test: %v", err)
	}

	// Fetch the user
	fetchedUser, err := store.GetUser(user.UserID)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	if fetchedUser.UserID != user.UserID || fetchedUser.Email != user.Email {
		t.Fatalf("Fetched user does not match created user. Expected user ID: %v, got: %v. Expected email: %v, got: %v", user.UserID, fetchedUser.UserID, user.Email, fetchedUser.Email)
	}

	// Cleanup
	err = store.DeleteUser(user.UserID)
	if err != nil {
		t.Fatalf("Failed to clean up test user: %v", err)
	}
}

func TestDeleteUser(t *testing.T) {
	db := testutils.SetupDB()
	defer db.Close()

	store := &Store{db: db}

	// Create a user first
	user := User{
		UserID:           "deleteTestUserID",
		ExternalID:       "deleteTestExternalID",
		Username:         "deleteTestUsername",
		PasswordEnabled:  true,
		Email:            "deleteTest@example.com",
		FirstName:        "DeleteTest",
		LastName:         "User",
		ImageURL:         "https://deleteTest.com/image.jpg",
		ProfileImageURL:  "https://deleteTest.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err := store.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user for delete test: %v", err)
	}

	// Delete the user
	err = store.DeleteUser(user.UserID)
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	// Fetch the user and check if it exists
	_, err = store.GetUser(user.UserID)
	if err == nil {
		t.Fatalf("User deletion failed. User still exists")
	}
}
