// //go:build integration
// // +build integration
package userservice

import (
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"testing"
)

func setupTests() (UserService, storemodels.User, error, *postgresqlclient.PostgreSQL) {
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

	db, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		return nil, user, err, db
	}
	userService, err := NewUserService()
	if err != nil {
		return nil, user, err, db
	}
	return userService, user, nil, db
}

func cleanupTests(user storemodels.User, db *postgresqlclient.PostgreSQL) error {
	return db.DeleteUserInUsersTable(user.UserID)
}

func TestCreateUser(t *testing.T) {
	userService, user, err, db := setupTests()
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
	userService, user, err, db := setupTests()
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

func TestGetUser(t *testing.T) {
	userService, user, err, db := setupTests()
	if err != nil {
		t.Error(err)
	}
	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	err = db.CheckIfUserInUsersTable(user.UserID, 5)
	if err != nil {
		t.Error(err)
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

func TestDeleteUser(t *testing.T) {
	userService, user, err, db := setupTests()
	if err != nil {
		t.Error(err)
	}
	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
	err = db.CheckIfUserInUsersTable(user.UserID, 5)
	if err != nil {
		t.Error(err)
	}
	err = userService.DeleteUser(user.UserID)
	if err != nil {
		t.Error(err)
	}
	err = db.CheckUserDeletedInUsersTable(user.UserID, 5)
	if err != nil {
		t.Error(err)
	}

	t.Cleanup(func() {
		cleanupTests(user, db)
	})
}
