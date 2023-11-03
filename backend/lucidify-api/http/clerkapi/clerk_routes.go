package clerkapi

import (
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/server/config"
	"lucidify-api/server/middleware"
	"lucidify-api/service/userservice"
	"net/http"
)

func SetupRoutes(storeInstance *postgresqlclient.PostgreSQL, userService userservice.UserService, config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	handler := ClerkHandler(storeInstance, userService)

	handler = middleware.ClerkWebhooksAuthenticationMiddleware(config)(handler)
	handler = middleware.Logging(handler)

	mux.HandleFunc("/clerk/webhook", handler)

	return mux
}
