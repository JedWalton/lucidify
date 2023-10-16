// //go:build integration
// // +build integration
package userservice

import (
	"lucidify-api/modules/store/postgresqlclient"
	"testing"
)

func setupTests() (UserService, error) {
	db, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		return nil, err
	}
	userService := NewUserService(db)
	return userService, nil
}

func TestCreateUser(t *testing.T) {
	userService, err := setupTests()
	if err != nil {
		t.Error(err)
	}
	err = userService.CreateUser()
	if err != nil {
		t.Error(err)
	}

}
