package server

import (
	"lucidify-api/api/chat"
	"lucidify-api/api/clerk"
	"lucidify-api/api/documents"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"
)

func SetupRoutes(storeInstance *store.Store, config *config.ServerConfig, mux *http.ServeMux) {
	chat.SetupRoutes(config, mux)
	documents.SetupRoutes(storeInstance, config, mux)
	clerk.SetupRoutes(storeInstance, config, mux)
}
