package userservice

import (
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/storemodels"
)

type UserService interface {
	CreateUser(user storemodels.User) error
	UpdateUser(user storemodels.User) error
	DeleteUser(userID string) error
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
	err := u.postgresqlDB.DeleteUserInUsersTable(userID)
	if err != nil {
		return err
	}
	return nil
}
