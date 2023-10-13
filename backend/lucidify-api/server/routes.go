package server

import (
	"lucidify-api/api/chatapi"
	"lucidify-api/api/clerkapi"
	"lucidify-api/api/documentsapi"
	"lucidify-api/modules/chatservice"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/store"
	"lucidify-api/modules/store/weaviateclient"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	storeInstance *postgresqlclient.PostgreSQL,
	clerkInstance clerk.Client,
	weaviateInstance weaviateclient.WeaviateClient,
	documentsService store.DocumentService,
	chatService chatservice.ChatService) {

	chatapi.SetupRoutes(config, mux, chatService, clerkInstance)
	documentsapi.SetupRoutes(config, mux, documentsService, clerkInstance)
	clerkapi.SetupRoutes(storeInstance, config, mux)
}
