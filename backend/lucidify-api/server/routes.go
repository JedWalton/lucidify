package server

import (
	"lucidify-api/api/chat"
	"lucidify-api/api/clerkapi"
	"lucidify-api/api/documentsapi"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/store"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	storeInstance *postgresqlclient.PostgreSQL,
	clerkInstance clerk.Client,
	documentsService store.DocumentService) {

	chat.SetupRoutes(config, mux)
	documentsapi.SetupRoutes(config, mux, documentsService, clerkInstance)
	clerkapi.SetupRoutes(storeInstance, config, mux)
}
