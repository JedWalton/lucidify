package server

import (
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/server/config"
	"lucidify-api/service/chatservice"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
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

	clerkInstance, err := clerkservice.NewClerkClient()
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

	SetupRoutes(
		config,
		mux,
		postgresqlDB,
		clerkInstance.GetClerkClient(),
		documentService,
		chatService,
	)

	BasicLogging(config, mux)
}
