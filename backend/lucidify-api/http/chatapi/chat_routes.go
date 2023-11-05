package chatapi

import (
	"lucidify-api/server/config"
	"lucidify-api/server/middleware"
	"lucidify-api/service/chatservice"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	cvs chatservice.ChatVectorService,
	clerkInstance clerk.Client) *http.ServeMux {

	mux = SetupChatHandler(config, mux, cvs, clerkInstance)

	return mux
}

// ... other code ...

func SetupChatHandler(
	config *config.ServerConfig,
	mux *http.ServeMux,
	cvs chatservice.ChatVectorService,
	clerkInstance clerk.Client) *http.ServeMux {

	handler := ChatHandler(clerkInstance, cvs)

	injectActiveSession := clerk.WithSession(clerkInstance)

	handler = middleware.Logging(handler)

	mux.Handle("/chat", injectActiveSession(handler))
	// mux.Handle("/api/sync", handler)

	return mux
}
