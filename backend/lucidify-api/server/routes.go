package server

import (
	"lucidify-api/api/chat"
	"lucidify-api/api/clerkapi"
	"lucidify-api/api/documents"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, storeInstance *store.Store, clerkInstance *clerk.Client) {
	chat.SetupRoutes(config, mux)
	documents.SetupRoutes(config, mux, storeInstance, clerkInstance)
	clerkapi.SetupRoutes(storeInstance, config, mux)
}
