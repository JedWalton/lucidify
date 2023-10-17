// //go:build integration
// // +build integration
package userservice

import (
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/storemodels"
	"testing"
)

func setupTests() (UserService, storemodels.User, error) {
	user := storemodels.User{
		UserID:           "TestCreateUserTableUserID",
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
		return nil, user, err
	}
	userService := NewUserService(db)
	return userService, user, nil
}

func TestCreateUser(t *testing.T) {
	userService, user, err := setupTests()
	if err != nil {
		t.Error(err)
	}
	err = userService.CreateUser(user)
	if err != nil {
		t.Error(err)
	}
}
