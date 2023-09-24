package documents

import (
	"lucidify-api/modules/clerkclient"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, storeInstance *store.Store) *http.ServeMux {

	client, err := clerkclient.NewClerkClient(config.ClerkSecretKey)
	if err != nil {
		panic(err)
	}
	// handler := DocumentsUploadHandler(storeInstance, client)

	injectActiveSession := clerk.WithSession(client)

	// // Wrap the handler with other middlewares
	// handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	// handler = middleware.Logging(handler)

	mux.Handle("/documents/upload", injectActiveSession(DocumentsUploadHandler(storeInstance, client)))

	return mux
}
