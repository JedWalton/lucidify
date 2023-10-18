package server

import (
	"log"
	"lucidify-api/server/config"
	"net/http"
)

func BasicLogging(config *config.ServerConfig, mux *http.ServeMux) {
	log.Printf("Server starting on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
