package server

import (
	"log"
	"lucidify-api/modules/chatservice"
	"lucidify-api/modules/clerkclient"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/store"
	"lucidify-api/modules/store/weaviateclient"
	"net/http"
)

func StartServer() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	postgresqlDB, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	clerkInstance, err := clerkclient.NewClerkClient(config.ClerkSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	weaviateInstance, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		log.Fatal(err)
	}

	documentsService := store.NewDocumentService(postgresqlDB, weaviateInstance)

	chatService := chatservice.NewChatService(postgresqlDB, weaviateInstance)

	SetupRoutes(
		config,
		mux,
		postgresqlDB,
		clerkInstance,
		weaviateInstance,
		documentsService,
		chatService,
	)

	BasicLogging(config, mux)
}
