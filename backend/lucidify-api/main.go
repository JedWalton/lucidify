package main

import (
	"log"
	"lucidify-api/config"
	"lucidify-api/lucidifychat"
	"lucidify-api/store"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func NewServerConfig() *config.ServerConfig {

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost",
	}

	Store := store.NewStore()

	return &config.ServerConfig{
		OPENAI_API_KEY: OPENAI_API_KEY,
		AllowedOrigins: allowedOrigins,
		Port:           port,
		Store:          Store,
	}
}

func main() {
	config := NewServerConfig()

	mux := http.NewServeMux()
	lucidifychat.SetupRoutes(config, mux)

	log.Printf("Server starting on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
