package chat

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/chat", middleware.Chain(
		ChatHandler(),
		middleware.CORSMiddleware(config.AllowedOrigins),
		middleware.Logging,
	))

	return mux
}
