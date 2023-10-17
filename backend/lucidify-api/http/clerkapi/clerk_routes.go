package clerkapi

import (
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/server/config"
	middleware2 "lucidify-api/server/middleware"
	"lucidify-api/service/userservice"
	"net/http"
)

func SetupRoutes(storeInstance *postgresqlclient.PostgreSQL, config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	userService, err := userservice.NewUserService()
	if err != nil {
		log.Fatal(err)
	}
	handler := ClerkHandler(storeInstance, userService)

	handler = middleware2.ClerkWebhooksAuthenticationMiddleware(config)(handler)
	handler = middleware2.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware2.Logging(handler)

	mux.HandleFunc("/clerk/webhook", handler)

	return mux
}
