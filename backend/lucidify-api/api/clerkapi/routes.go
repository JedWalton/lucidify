package clerkapi

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store/postgresqlclient"
	"net/http"
)

func SetupRoutes(storeInstance *postgresqlclient.PostgreSQL, config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	handler := ClerkHandler(storeInstance)

	handler = middleware.ClerkWebhooksAuthenticationMiddleware(config)(handler)
	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.HandleFunc("/clerk/webhook", handler)

	return mux
}
