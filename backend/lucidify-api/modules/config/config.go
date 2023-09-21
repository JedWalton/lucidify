package config

import (
	"log"
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type ServerConfig struct {
	OPENAI_API_KEY     string
	AllowedOrigins     []string
	Port               string
	PostgresqlURL      string
	ClerkClient        clerk.Client
	ClerkSigningSecret string
	ClerkSecretKey     string
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

	postgresqlURL := os.Getenv("POSTGRESQL_URL")
	if postgresqlURL == "" {
		log.Fatal("POSTGRESQL_URL environment variable is not set")
	}

	clerkClient, err := clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
	if err != nil {
		log.Fatalf("Failed to create Clerk client: %v", err)
	}

	clerkSigningSecret := os.Getenv("CLERK_SIGNING_SECRET")
	if clerkSigningSecret == "" {
		log.Fatal("CLERK_SIGNING_SECRET environment variable is not set")
	}

	clerkSecretKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecretKey == "" {
		log.Fatalf("CLERK_SECRET_KEY environment variable is not set: %v", err)
	}

	return &ServerConfig{
		OPENAI_API_KEY:     OPENAI_API_KEY,
		AllowedOrigins:     allowedOrigins,
		Port:               port,
		PostgresqlURL:      postgresqlURL,
		ClerkClient:        clerkClient,
		ClerkSigningSecret: clerkSigningSecret,
		ClerkSecretKey:     clerkSecretKey,
	}
}
