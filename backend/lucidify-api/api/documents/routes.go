package documents

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store"
	"net/http"
)

func SetupRoutes(storeInstance *store.Store, config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/documents/upload", middleware.Chain(
		DocumentsUploadHandler(storeInstance),
		middleware.CORSMiddleware(config.AllowedOrigins),
		middleware.Logging,
	))

	return mux
}
