package main

import (
	"log"
	"lucidify-api/api/chat"
	"lucidify-api/api/clerk"
	"lucidify-api/api/documents"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	store, err := store.NewStore(os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	SetupRoutes(config, mux, store)

	log.Printf("Server starting on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, store *store.Store) {
	chat.SetupRoutes(config, mux, store)
	documents.SetupRoutes(config, mux, store)
	clerk.SetupRoutes(config, mux, store)
}
