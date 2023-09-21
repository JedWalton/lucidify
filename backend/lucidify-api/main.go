package main

import (
	"log"
	"lucidify-api/api/chat"
	"lucidify-api/api/clerk"
	"lucidify-api/api/documents"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"
)

func main() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	storeInstance, err := store.NewStore(config.PostgresqlURL)
	if err != nil {
		log.Fatal(err)
	}
	SetupRoutes(storeInstance, config, mux)

	log.Printf("Server starting on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func SetupRoutes(storeInstance *store.Store, config *config.ServerConfig, mux *http.ServeMux) {
	chat.SetupRoutes(config, mux)
	documents.SetupRoutes(storeInstance, config, mux)
	clerk.SetupRoutes(storeInstance, config, mux)
}
