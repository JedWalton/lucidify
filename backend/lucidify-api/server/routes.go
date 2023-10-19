package server

import (
	"lucidify-api/http/chatapi"
	"lucidify-api/http/clerkapi"
	"lucidify-api/http/documentsapi"
	"lucidify-api/server/config"
	"lucidify-api/service/chatservice"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
	"lucidify-api/service/userservice"
	"net/http"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	clerkService clerkservice.ClerkClient,
	documentsService documentservice.DocumentService,
	chatService chatservice.ChatService,
	userService userservice.UserService) {

	chatapi.SetupRoutes(config, mux, chatService, clerkService.GetClerkClient())
	documentsapi.SetupRoutes(config, mux, documentsService, clerkService)
	clerkapi.SetupRoutes(config, mux, userService)
}
