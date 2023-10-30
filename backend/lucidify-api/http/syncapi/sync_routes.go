package syncapi

import (
	"lucidify-api/server/config"
	"lucidify-api/server/middleware"
	"lucidify-api/service/syncservice"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	clerkInstance clerk.Client,
	syncService syncservice.SyncService) *http.ServeMux {

	mux = SetupSyncHandler(config, mux, syncService)

	return mux
}

func SetupSyncHandler(config *config.ServerConfig,
	mux *http.ServeMux, syncService syncservice.SyncService) *http.ServeMux {

	handler := SyncHandler(syncService)

	handler = middleware.Logging(handler)

	// mux.Handle("/api/sync", handler)
	mux.Handle("/api/sync/localstorage/", http.StripPrefix("/api/sync/localstorage/", handler))

	return mux
}
