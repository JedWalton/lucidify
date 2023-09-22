package documents

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, storeInstance *store.Store, clerkInstance *clerk.Client) *http.ServeMux {
	handler := DocumentsUploadHandler(storeInstance)

	// Wrap the handler with Clerk's authentication middleware
	handler = func(w http.ResponseWriter, r *http.Request) {
		clerk.RequireSessionV2(*clerkInstance)(http.HandlerFunc(handler)).ServeHTTP(w, r)
	}

	// Wrap the handler with other middlewares
	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.HandleFunc("/documents/upload", handler)

	return mux
}
