package server

import (
	"log"
	"lucidify-api/modules/clerkclient"
	"lucidify-api/modules/config"
	"lucidify-api/modules/postgresqlclient"
	"net/http"
)

func StartServer() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	storeInstance, err := postgresqlclient.NewPostgreSQL(config.PostgresqlURL)
	if err != nil {
		log.Fatal(err)
	}

	clerkInstance, err := clerkclient.NewClerkClient(config.ClerkSecretKey)
	if err != nil {
		log.Fatal(err)
	}

	// weaviateInstance, err := weaviateclient.NewWeaviateClient()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	SetupRoutes(config, mux, storeInstance, clerkInstance)

	BasicLogging(config, mux)
}
