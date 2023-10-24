package syncapi

import (
	"lucidify-api/server/config"
	"lucidify-api/server/middleware"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	clerkInstance clerk.Client) *http.ServeMux {

	mux = SetupSyncHandler(config, mux)

	return mux
}

func SetupSyncHandler(config *config.ServerConfig, mux *http.ServeMux) *http.ServeMux {

	handler := SyncHandler()

	handler = middleware.Logging(handler)

	// mux.Handle("/api/sync", handler)
	mux.Handle("/api/sync/", http.StripPrefix("/api/sync/", handler))

	return mux
}
