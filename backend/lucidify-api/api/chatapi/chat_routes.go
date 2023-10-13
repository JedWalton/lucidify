package chatapi

import (
	"lucidify-api/modules/chatservice"
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	chatService chatservice.ChatService,
	clerkInstance clerk.Client) *http.ServeMux {

	mux = SetupChatHandler(config, mux, chatService, clerkInstance)

	return mux
}

func SetupChatHandler(
	config *config.ServerConfig,
	mux *http.ServeMux,
	chatService chatservice.ChatService,
	clerkInstance clerk.Client) *http.ServeMux {

	handler := ChatHandler(clerkInstance, chatService)

	injectActiveSession := clerk.WithSession(clerkInstance)

	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.Handle("/chat", injectActiveSession(handler))

	return mux
}
