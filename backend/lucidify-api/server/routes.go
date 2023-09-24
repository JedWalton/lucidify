package server

import (
	"lucidify-api/api/chat"
	"lucidify-api/api/clerkapi"
	"lucidify-api/api/documents"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, storeInstance *store.Store) {
	chat.SetupRoutes(config, mux)
	documents.SetupRoutes(config, mux, storeInstance)
	clerkapi.SetupRoutes(storeInstance, config, mux)
}
