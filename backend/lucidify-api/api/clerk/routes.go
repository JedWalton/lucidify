package clerk

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/clerk/webhook", middleware.Chain(
		ClerkHandler(config.Store),
		middleware.CORSMiddleware(config.AllowedOrigins),
		middleware.Logging,
	))

	return mux
}
