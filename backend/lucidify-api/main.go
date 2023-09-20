package main

import (
	"log"
	"lucidify-api/api/chat"
	"lucidify-api/api/clerk"
	"lucidify-api/api/documents"
	"lucidify-api/modules/config"
	"net/http"
)

func main() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	SetupRoutes(config, mux)

	log.Printf("Server starting on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux) {
	chat.SetupRoutes(config, mux)
	documents.SetupRoutes(config, mux, config.Store)
	clerk.SetupRoutes(config, mux)
}
