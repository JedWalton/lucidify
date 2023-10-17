package userservice

import (
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/storemodels"
)

type UserService interface {
	CreateUser(user storemodels.User) error
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
