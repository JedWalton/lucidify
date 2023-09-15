package main

import (
	"log"
	"lucidify-api/api/chat"
	"lucidify-api/api/documents"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	store := store.NewStore()

	SetupRoutes(config, mux, store)

	log.Printf("Server starting on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, store *store.Store) {
	chat.SetupRoutes(config, mux, store)
	documents.SetupRoutes(config, mux, store)
}
