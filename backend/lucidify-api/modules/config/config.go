package config

import (
	"log"
	"lucidify-api/modules/store"
	"os"
)

type ServerConfig struct {
	OPENAI_API_KEY string
	AllowedOrigins []string
	Port           string
	Store          *store.Store
}

func NewServerConfig() *ServerConfig {

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
	}

	Store := store.NewStore()

	return &ServerConfig{
		OPENAI_API_KEY: OPENAI_API_KEY,
		AllowedOrigins: allowedOrigins,
		Port:           port,
		Store:          Store,
	}
}
