package chat

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, store *store.Store) *http.ServeMux {
	mux.HandleFunc("/chat", middleware.Chain(
		ChatHandler(),
		middleware.CORSMiddleware(config.AllowedOrigins),
		middleware.Logging,
	))

	return mux
}
