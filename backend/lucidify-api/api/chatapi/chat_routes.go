package chatapi

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store/weaviateclient"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	weaviateInstance weaviateclient.WeaviateClient,
	clerkInstance clerk.Client) *http.ServeMux {

	mux = SetupChatHandler(config, mux, weaviateInstance, clerkInstance)

	return mux
}

func SetupChatHandler(
	config *config.ServerConfig,
	mux *http.ServeMux,
	weaviateInstance weaviateclient.WeaviateClient,
	clerkInstance clerk.Client) *http.ServeMux {

	handler := ChatHandler(clerkInstance)

	injectActiveSession := clerk.WithSession(clerkInstance)

	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.Handle("/chat", injectActiveSession(handler))

	return mux
}
