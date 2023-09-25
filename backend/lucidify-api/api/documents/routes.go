package documents

import (
	"lucidify-api/modules/config"
	"lucidify-api/modules/middleware"
	"lucidify-api/modules/store"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, storeInstance *store.Store, client clerk.Client) *http.ServeMux {

	handler := DocumentsUploadHandler(storeInstance, client)

	injectActiveSession := clerk.WithSession(client)

	// // Wrap the handler with other middlewares
	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.Handle("/documents/upload", injectActiveSession(handler))

	return mux
}
