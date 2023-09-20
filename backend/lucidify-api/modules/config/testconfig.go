package config

import (
	"log"
	"lucidify-api/modules/store"
	"os"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/joho/godotenv"
)

type TestServerConfig struct {
	OPENAI_API_KEY string
	AllowedOrigins []string
	Port           string
	TestStore      *store.Store
	ClerkClient    clerk.Client
}

func NewTestServerConfig() *TestServerConfig {
	// Load environment variables from the .env file
	if err := godotenv.Load("../../../../.env"); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

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

	testStore, err := store.NewStore(postgresqlURL)
	if err != nil {
		log.Fatal(err)
	}

	clerkClient, err := clerk.NewClient(os.Getenv("CLERK_SECRET_KEY"))
	if err != nil {
		log.Fatalf("Failed to create Clerk client: %v", err)
	}

	return &TestServerConfig{
		OPENAI_API_KEY: OPENAI_API_KEY,
		AllowedOrigins: allowedOrigins,
		Port:           port,
		TestStore:      testStore,
		ClerkClient:    clerkClient,
	}
}
