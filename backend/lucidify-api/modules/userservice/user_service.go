package userservice

import (
	"fmt"
	"lucidify-api/modules/store/postgresqlclient"
)

type UserService interface {
	CreateUser() error
}

type UserServiceImpl struct {
	postgresqlDB *postgresqlclient.PostgreSQL
}

func NewUserService(
	postgresqlDB *postgresqlclient.PostgreSQL) UserService {
	return &UserServiceImpl{postgresqlDB: postgresqlDB}
}

func (u *UserServiceImpl) CreateUser() error {
	return fmt.Errorf("not implemented")
}
