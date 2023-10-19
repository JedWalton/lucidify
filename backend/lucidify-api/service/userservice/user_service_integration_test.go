// //go:build integration
// // +build integration
package userservice

import (
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/service/documentservice"
	"testing"
	"time"
)

func setupTests() (UserService, storemodels.User, error, *postgresqlclient.PostgreSQL, documentservice.DocumentService) {
	user := storemodels.User{
		UserID:           "TestCreateUserTableUserServiceUserID",
		ExternalID:       "TestCreateUserTableExternalID",
		Username:         "TestCreateUserableUsername",
		PasswordEnabled:  true,
		Email:            "TestCreateUser@example.com",
		FirstName:        "TestCreateUserCreateTest",
		LastName:         "TestCreateUserUser",
		ImageURL:         "https://TestCreateUser.com/image.jpg",
		ProfileImageURL:  "https://TestCreateUser.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		return nil, user, err, postgresqlDB, nil
	}

	weaviateDB, err := weaviateclient.NewWeaviateClientTest()
	docService := documentservice.NewDocumentService(postgresqlDB, weaviateDB)

	userService := NewUserService(postgresqlDB)
	if err != nil {
		return nil, user, err, postgresqlDB, nil
	}

	userService.SetDocumentService(docService)
	return userService, user, nil, postgresqlDB, docService
}

func cleanupTests(user storemodels.User, db *postgresqlclient.PostgreSQL) error {
	return db.DeleteUserInUsersTable(user.UserID)
}

func TestCreateUser(t *testing.T) {
	userService, user, err, db, _ := setupTests()
	if err != nil {
		t.Error(err)
	}
	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	userFromDb, err := db.GetUserInUsersTable(user.UserID)
	if err != nil {
		t.Error(err)
	}
	if userFromDb.UserID != user.UserID {
		t.Errorf("Expected %s, got %s", user.UserID, userFromDb.UserID)
	}

	t.Cleanup(func() {
		cleanupTests(user, db)
	})
}

func TestUpdateUser(t *testing.T) {
	userService, user, err, db, _ := setupTests()
	if err != nil {
		t.Error(err)
	}
	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	userUpdated := storemodels.User{
		UserID:           "TestCreateUserTableUserServiceUserID",
		ExternalID:       "TestCreateUserTableExternalIDUpdated",
		Username:         "TestCreateUserableUsernameUpdated",
		PasswordEnabled:  true,
		Email:            "TestCreateUserUpdated@example.com",
		FirstName:        "TestCreateUserCreateTestUpdated",
		LastName:         "TestCreateUserUserUpdated",
		ImageURL:         "https://TestCreateUserUpdated.com/image.jpg",
		ProfileImageURL:  "https://TestCreateUserUpdated.com/profile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}
	err = userService.UpdateUser(userUpdated)
	if err != nil {
		t.Error(err)
	}
	userAfterUpdate, err := db.GetUserInUsersTable(user.UserID)
	if err != nil {
		t.Error(err)
	}
	if userAfterUpdate.UserID != user.UserID {
		t.Errorf("Expected %s, got %s", user.UserID, userAfterUpdate.UserID)
	}
	if userAfterUpdate.ExternalID != userUpdated.ExternalID {
		t.Errorf("Expected %s, got %s", userUpdated.ExternalID, userAfterUpdate.ExternalID)
	}
	if userAfterUpdate.Username != userUpdated.Username {
		t.Errorf("Expected %s, got %s", userUpdated.Username, userAfterUpdate.Username)
	}
	if userAfterUpdate.PasswordEnabled != userUpdated.PasswordEnabled {
		t.Errorf("Expected %t, got %t", userUpdated.PasswordEnabled, userAfterUpdate.PasswordEnabled)
	}
	if userAfterUpdate.Email != userUpdated.Email {
		t.Errorf("Expected %s, got %s", userUpdated.Email, userAfterUpdate.Email)
	}
	if userAfterUpdate.FirstName != userUpdated.FirstName {
		t.Errorf("Expected %s, got %s", userUpdated.FirstName, userAfterUpdate.FirstName)
	}
	if userAfterUpdate.LastName != userUpdated.LastName {
		t.Errorf("Expected %s, got %s", userUpdated.LastName, userAfterUpdate.LastName)
	}
	if userAfterUpdate.ImageURL != userUpdated.ImageURL {
		t.Errorf("Expected %s, got %s", userUpdated.ImageURL, userAfterUpdate.ImageURL)
	}
	if userAfterUpdate.ProfileImageURL != userUpdated.ProfileImageURL {
		t.Errorf("Expected %s, got %s", userUpdated.ProfileImageURL, userAfterUpdate.ProfileImageURL)
	}
	if userAfterUpdate.TwoFactorEnabled != userUpdated.TwoFactorEnabled {
		t.Errorf("Expected %t, got %t", userUpdated.TwoFactorEnabled, userAfterUpdate.TwoFactorEnabled)
	}
	if userAfterUpdate.CreatedAt != userUpdated.CreatedAt {
		t.Errorf("Expected %d, got %d", userUpdated.CreatedAt, userAfterUpdate.CreatedAt)
	}
	if userAfterUpdate.UpdatedAt != userUpdated.UpdatedAt {
		t.Errorf("Expected %d, got %d", userUpdated.UpdatedAt, userAfterUpdate.UpdatedAt)
	}

	t.Cleanup(func() {
		cleanupTests(user, db)
	})
}

func TestDeleteUser(t *testing.T) {
	userService, user, err, postgresqlDB, _ := setupTests()
	if err != nil {
		t.Error(err)
	}

	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	_, err = userService.GetUserWithRetries(user.UserID, 5)
	if err != nil {
		t.Errorf("User not found after creation: %v", err)
	}
	err = userService.DeleteUser(user.UserID)
	if err != nil {
		t.Error(err)
	}
	if !userService.HasUserBeenDeleted(user.UserID, 5) {
		t.Errorf("User not deleted after deletion: %v", err)
	}

	// Verify All Documents, Chunks in Postgresql and Weaviate are deleted

	t.Cleanup(func() {
		cleanupTests(user, postgresqlDB)
	})
}

func TestGetUser(t *testing.T) {
	userService, user, err, db, _ := setupTests()
	if err != nil {
		t.Error(err)
	}
	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	_, err = userService.GetUserWithRetries(user.UserID, 5)
	if err != nil {
		t.Errorf("User not found after creation: %v", err)
	}
	userFromDb, err := userService.GetUser(user.UserID)
	if err != nil {
		t.Error(err)
	}
	if userFromDb.Email != user.Email {
		t.Errorf("Expected %s, got %s", user.Email, userFromDb.Email)
	}

	t.Cleanup(func() {
		cleanupTests(user, db)
	})
}

func TestGetUserWithRetries(t *testing.T) {
	userService, user, err, db, _ := setupTests()
	if err != nil {
		t.Error(err)
	}

	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	// Case 1: User is found without needing retries
	foundUser, err := userService.GetUserWithRetries(user.UserID, 5)
	if err != nil || foundUser == nil {
		t.Error("GetUserWithRetries failed to retrieve the user, but the user should have been found")
	}
	// Case 2: User is not found even after retries
	nonExistentUserID := "non-existent-user-id"
	// Delay the start of the next case to avoid immediate execution after user creation
	time.Sleep(2 * time.Second)

	_, err = userService.GetUserWithRetries(nonExistentUserID, 5)
	if err == nil {
		t.Error("GetUserWithRetries did not return an error, but it should have since the user does not exist")
	}

	t.Cleanup(func() {
		cleanupTests(user, db)
	})
}

func TestHasUserBeenDeleted(t *testing.T) {
	userService, user, err, db, _ := setupTests()
	if err != nil {
		t.Error(err)
	}

	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}

	// Case 1: User has not been deleted
	userNotDeleted := userService.HasUserBeenDeleted(user.UserID, 5)
	if userNotDeleted {
		t.Error("HasUserBeenDeleted returns true, but the user was not deleted")
	}

	// Assume we have a function to delete a user for test
	userService.DeleteUser(user.UserID)

	// Case 2: User has been deleted
	userHasBeenDeleted := userService.HasUserBeenDeleted(user.UserID, 5)
	if !userHasBeenDeleted {
		t.Error("HasUserBeenDeleted returns false, but the user was deleted")
	}

	t.Cleanup(func() {
		cleanupTests(user, db)
	})
}
