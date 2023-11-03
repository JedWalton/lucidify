package userservice

import (
	"fmt"
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/storemodels"
	"lucidify-api/data/store/weaviateclient"
	"time"

	"github.com/google/uuid"
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
	weaviateDB   weaviateclient.WeaviateClient
}

func NewUserService(weaviateClient weaviateclient.WeaviateClient) (UserService, error) {
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

func (u *UserServiceImpl) deleteDocument(documentID uuid.UUID) error {
	chunks, err := u.postgresqlDB.GetChunksOfDocumentByDocumentID(documentID)
	if err != nil {
		return fmt.Errorf("Failed to get chunks of document: %w", err)
	}
	err = u.weaviateDB.DeleteChunks(chunks)
	if err != nil {
		return fmt.Errorf("Failed to delete chunks from Weaviate: %w", err)
	}
	err = u.postgresqlDB.DeleteDocumentByUUID(documentID)
	if err != nil {
		log.Printf("Failed to delete document from PostgreSQL: %v", err)
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(userID string) error {
	documents, err := u.postgresqlDB.GetAllDocuments(userID)
	if err != nil {
		return fmt.Errorf("Failed to get all documents from PostgreSQL: %w", err)
	}
	for _, document := range documents {
		if err := u.deleteDocument(document.DocumentUUID); err != nil {
			return fmt.Errorf("Failed to delete document: %w", err)
		}
	}

	err = u.postgresqlDB.DeleteUserInUsersTable(userID)
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
