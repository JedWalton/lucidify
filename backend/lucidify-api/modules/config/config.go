package config

import (
	"log"
	"lucidify-api/modules/store"
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type ServerConfig struct {
	OPENAI_API_KEY string
	AllowedOrigins []string
	Port           string
	Store          *store.Store
	ClerkClient    clerk.Client
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

	store, err := store.NewStore(os.Getenv("POSTGRESQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	clerkClient, err := clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
	if err != nil {
		log.Fatalf("Failed to create Clerk client: %v", err)
	}

	return &ServerConfig{
		OPENAI_API_KEY: OPENAI_API_KEY,
		AllowedOrigins: allowedOrigins,
		Port:           port,
		Store:          store,
		ClerkClient:    clerkClient,
	}
}
