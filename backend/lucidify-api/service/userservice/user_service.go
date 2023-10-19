package userservice

import (
	"fmt"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"time"
)

type UserService interface {
	CreateUser(user storemodels.User) error
	UpdateUser(user storemodels.User) error
	DeleteUser(userID string) error
	GetUser(userID string) (*storemodels.User, error)
	GetUserWithRetries(userID string, retries int) (*storemodels.User, error)
	HasUserBeenDeleted(userID string, retries int) bool
}

type UserServiceImpl struct {
	postgresqlDB *postgresqlclient.PostgreSQL
}

func NewUserService() (UserService, error) {
	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		return nil, err
	}
	return &UserServiceImpl{postgresqlDB: postgresqlDB}, nil
}

func (u *UserServiceImpl) CreateUser(user storemodels.User) error {
	err := u.postgresqlDB.CreateUserInUsersTable(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServiceImpl) UpdateUser(user storemodels.User) error {
	err := u.postgresqlDB.UpdateUserInUsersTable(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(userID string) error {
	// Perform Full Cleanup Of All User Data
	// Weaviate Data
	// Postgres Data
	// Confirm Cleanup
	// Then Delete User from Users Table
	err := u.postgresqlDB.DeleteUserInUsersTable(userID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserServiceImpl) GetUser(userID string) (*storemodels.User, error) {
	user, err := u.postgresqlDB.GetUserInUsersTable(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserServiceImpl) GetUserWithRetries(userID string, retries int) (*storemodels.User, error) {
	var found bool
	var user *storemodels.User
	var err error

	for i := 0; i < retries; i++ {
		user, err = u.GetUser(userID)
		if err == nil {
			found = true
			break
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	if found {
		return user, nil
	}
	return nil, fmt.Errorf("User not found after %d retries", retries)
}

func (u *UserServiceImpl) HasUserBeenDeleted(userID string, retries int) bool {
	for i := 0; i < retries; i++ {
		_, err := u.GetUser(userID)
		if err != nil {
			return true
		}
	}
	return false
}
