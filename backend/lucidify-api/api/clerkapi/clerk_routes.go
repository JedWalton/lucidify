package clerkapi

import (
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/userservice"
	"net/http"
)

func SetupRoutes(storeInstance *postgresqlclient.PostgreSQL, config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	userService, err := userservice.NewUserService()
	if err != nil {
		log.Fatal(err)
	}
	handler := ClerkHandler(storeInstance, userService)

	handler = middleware.ClerkWebhooksAuthenticationMiddleware(config)(handler)
	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.HandleFunc("/clerk/webhook", handler)

	return mux
}
