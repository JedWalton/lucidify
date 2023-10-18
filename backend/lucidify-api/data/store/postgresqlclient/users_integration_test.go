// //go:build integration
// // +build integration
package postgresqlclient

import (
	"fmt"
	"lucidify-api/data/store/storemodels"
	"testing"
	"time"
)

func checkIfUserInUsersTable(userID string, retries int) error {
	store, err := NewPostgreSQL()
	if err != nil {
		return fmt.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	for i := 0; i < retries; i++ {
		_, err := store.GetUserInUsersTable(userID)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	return fmt.Errorf("User not found after %d retries", retries)
}

func TestCreateUserInUsersTable(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

	user := storemodels.User{
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
		t.Errorf("Failed to create user: %v", err)
	}

	// Check if the user exists
	err = checkIfUserInUsersTable(user.UserID, 3)
	if err != nil {
		t.Errorf("User not found after creation: %v", err)
	}

	// Register cleanup function
	t.Cleanup(func() {
		store.DeleteUserInUsersTable(user.UserID)
	})
}
func checkUserHasExpectedFirstNameAndLastNameInUsersTable(userID string, retries int, expectedFirstName string, expectedLastName string) error {
	db, err := NewPostgreSQL()
	if err != nil {
		return fmt.Errorf("Failed to create test postgresqlclient: %v", err)
	}
	for i := 0; i < retries; i++ {
		user, err := db.GetUserInUsersTable(userID)
		if err == nil && user.FirstName == expectedFirstName && user.LastName == expectedLastName {
			return nil
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	return fmt.Errorf("User not updated correctly after %d retries", retries)
}

func TestUpdateUserInUsersTable(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

	// Create a user first
	user := storemodels.User{
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
		t.Errorf("Failed to create user for update test: %v", err)
	}

	// Update the user
	user.FirstName = "UpdatedFirstName"
	user.LastName = "UpdatedLastName"
	err = store.UpdateUserInUsersTable(user)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}

	// Check if the user has the expected first name and last name
	err = checkUserHasExpectedFirstNameAndLastNameInUsersTable(user.UserID, 3, "UpdatedFirstName", "UpdatedLastName")
	if err != nil {
		t.Errorf("User not updated correctly: %v", err)
	}

	// Register cleanup function
	t.Cleanup(func() {
		store.DeleteUserInUsersTable(user.UserID)
	})

}

func TestGetUserInUsersTable(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

	// Create a user first
	user := storemodels.User{
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
		t.Errorf("Failed to create user for get test: %v", err)
	}

	// Fetch the user
	fetchedUser, err := store.GetUserInUsersTable(user.UserID)
	if err != nil {
		t.Errorf("Failed to fetch user: %v", err)
	}

	if fetchedUser.UserID != user.UserID || fetchedUser.Email != user.Email {
		t.Errorf("Fetched user does not match created user. Expected user ID: %v, got: %v. Expected email: %v, got: %v", user.UserID, fetchedUser.UserID, user.Email, fetchedUser.Email)
	}

	// Cleanup
	t.Cleanup(func() {
		store.DeleteUserInUsersTable(user.UserID)
	})
}

func TestDeleteUserInUsersTable(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

	// Create a user first
	user := storemodels.User{
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
		t.Errorf("Failed to create user for delete test: %v", err)
	}

	// Delete the user
	err = store.DeleteUserInUsersTable(user.UserID)
	if err != nil {
		t.Errorf("Failed to delete user: %v", err)
	}

	// Check if the user has been deleted
	var deleted bool
	for i := 0; i < 3; i++ {
		user, err := store.GetUserInUsersTable(user.UserID)
		if user == nil || err != nil {
			// If the user is not found, it means the user has been deleted
			deleted = true
			break
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	if !deleted {
		t.Errorf("User still exists in the database after 3 retries")
	}

	t.Cleanup(func() {
		store.DeleteUserInUsersTable(user.UserID)
	})
}
