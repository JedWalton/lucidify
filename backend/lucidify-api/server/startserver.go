package server

import (
	"log"
	"lucidify-api/modules/clerkclient"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"
)

func StartServer() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	storeInstance, err := store.NewStore(config.PostgresqlURL)
	if err != nil {
		log.Fatal(err)
	}

	clerkInstance, err := clerkclient.NewClerkClient(config.ClerkSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	SetupRoutes(config, mux, storeInstance, clerkInstance)

	BasicLogging(config, mux)
}
