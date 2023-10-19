package server

import (
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/server/config"
	"lucidify-api/service/chatservice"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
	"lucidify-api/service/userservice"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

func StartServer() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	clerkService, err := clerkservice.NewClerkClient()
	if err != nil {
		log.Fatal(err)
	}

	weaviateInstance, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		log.Fatal(err)
	}

	documentService := documentservice.NewDocumentService(postgresqlDB, weaviateInstance)

	openaiClient := openai.NewClient(config.OPENAI_API_KEY)

	chatService := chatservice.NewChatService(postgresqlDB, weaviateInstance, openaiClient, documentService)

	userService := userservice.NewUserService(postgresqlDB)
	userService.SetDocumentService(documentService)

	SetupRoutes(
		config,
		mux,
		clerkService,
		documentService,
		chatService,
		userService,
	)

	BasicLogging(config, mux)
}
