package documents

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, storeInstance *store.Store, clerkInstance *clerk.Client) *http.ServeMux {
	mux.HandleFunc("/documents/upload", middleware.Chain(
		DocumentsUploadHandler(storeInstance),
		middleware.CORSMiddleware(config.AllowedOrigins),
		middleware.Logging,
	))

	return mux
}
