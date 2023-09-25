package chat

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	handler := ChatHandler()

	// Wrap the handler with other middlewares
	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.HandleFunc("/chat", handler)

	return mux
}
