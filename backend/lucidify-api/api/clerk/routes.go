package clerk

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store"
	"net/http"
)

func SetupRoutes(storeInstance *store.Store, config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	handler := ClerkHandler(storeInstance)

	handler = middleware.ClerkWebhooksAuthenticationMiddleware(config)(handler)
	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.HandleFunc("/clerk/webhook", handler)

	return mux
}
