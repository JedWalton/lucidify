package server

import (
	"log"
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

	SetupRoutes(storeInstance, config, mux)

	BasicLogging(config, mux)
}

func StartTestServer() {
	config := config.NewTestServerConfig()

	mux := http.NewServeMux()

	storeInstance, err := store.NewStore(config.PostgresqlURL)
	if err != nil {
		log.Fatal(err)
	}

	SetupRoutes(storeInstance, config, mux)

	BasicLogging(config, mux)
}
