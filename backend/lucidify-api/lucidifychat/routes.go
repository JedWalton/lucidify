package lucidifychat

import (
	"lucidify-api/config"
	"lucidify-api/middleware"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/lucidifychat", middleware.Chain(
		LucidifyChatHandler(),
		middleware.CORSMiddleware(config.AllowedOrigins),
		middleware.Logging,
	))

	return mux
}
