package clerk

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	handler := ClerkHandler(config.Store)

	handler = middleware.ClerkWebhooksAuthenticationMiddleware(config)(handler)
	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.HandleFunc("/clerk/webhook", handler)

	return mux
}
