package chat

import (
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"
)

type ChatService interface {
	GenerateSystemPromptFromMessages() (string, error)
}

type ChatServiceImpl struct {
	postgresqlDB postgresqlclient.PostgreSQL
	weaviateDB   weaviateclient.WeaviateClient
}

func NewChatService(
	postgresqlDB *postgresqlclient.PostgreSQL,
	weaviateDB weaviateclient.WeaviateClient) ChatService {
	return &ChatServiceImpl{postgresqlDB: *postgresqlDB, weaviateDB: weaviateDB}
}

func (c *ChatServiceImpl) GenerateSystemPromptFromMessages() (string, error) {
	// do something
	return "chat", nil
}
