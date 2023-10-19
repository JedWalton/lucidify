package server

import (
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/http/chatapi"
	"lucidify-api/http/clerkapi"
	"lucidify-api/http/documentsapi"
	"lucidify-api/server/config"
	"lucidify-api/service/chatservice"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
	"net/http"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	storeInstance *postgresqlclient.PostgreSQL,
	clerkService clerkservice.ClerkClient,
	documentsService documentservice.DocumentService,
	chatService chatservice.ChatService) {

	chatapi.SetupRoutes(config, mux, chatService, clerkService.GetClerkClient())
	documentsapi.SetupRoutes(config, mux, documentsService, clerkService)
	clerkapi.SetupRoutes(storeInstance, config, mux)
}
