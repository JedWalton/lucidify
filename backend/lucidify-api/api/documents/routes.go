package documents

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/documents/upload", middleware.Chain(
		DocumentsUploadHandler(config.Store),
		middleware.CORSMiddleware(config.AllowedOrigins),
		middleware.Logging,
	))

	return mux
}
