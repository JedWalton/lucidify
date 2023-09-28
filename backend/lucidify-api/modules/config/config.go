package config

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	OPENAI_API_KEY      string
	AllowedOrigins      []string
	Port                string
	PostgresqlURL       string
	ClerkClient         clerk.Client
	ClerkSigningSecret  string
	ClerkSecretKey      string
	TestJWTSessionToken string
	TestUserID          string
}

func getGitRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func NewServerConfig() *ServerConfig {
	// Check if POSTGRESQL_URL environment variable is missing
	// if os.Getenv("POSTGRESQL_URL") == "" || os.Getenv("OPENAI_API_KEY") == "" || os.Getenv("CLERK_SECRET_KEY") == "" || os.Getenv("CLERK_SIGNING_SECRET") == "" {
	// 	// If missing, load the .env file
	// 	if err := godotenv.Load("../../../../.env"); err != nil {
	// 		log.Fatalf("Failed to load .env file: %v", err)
	// 	}
	// }
	if os.Getenv("POSTGRESQL_URL") == "" || os.Getenv("OPENAI_API_KEY") == "" || os.Getenv("CLERK_SECRET_KEY") == "" || os.Getenv("CLERK_SIGNING_SECRET") == "" {
		// If missing, load the .env file
		gitRoot, err := getGitRoot()
		if err != nil {
			log.Fatalf("Error getting git root: %v", err)
		}
		envPath := filepath.Join(gitRoot, ".env")
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("Failed to load .env file: %v", err)
		}
	}

	// Now retrieve the environment variables
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
		log.Fatalf("CLERK_SECRET_KEY environment variable is not set")
	}

	testJWTSessionToken := os.Getenv("TEST_JWT_SESSION_TOKEN")
	if testJWTSessionToken == "" {
		log.Printf("Either TEST_JWT_SESSION_TOKEN environment variable is not set " +
			"in development, or this is production and thus not an issue.")
	}

	testUserID := os.Getenv("TEST_USER_ID")
	if testUserID == "" {
		log.Printf("Either TEST_USER_ID environment variable is not set " +
			"in development, or this is production and thus not an issue.")
	}

	return &ServerConfig{
		OPENAI_API_KEY:      OPENAI_API_KEY,
		AllowedOrigins:      allowedOrigins,
		Port:                port,
		PostgresqlURL:       postgresqlURL,
		ClerkClient:         clerkClient,
		ClerkSigningSecret:  clerkSigningSecret,
		ClerkSecretKey:      clerkSecretKey,
		TestJWTSessionToken: testJWTSessionToken,
		TestUserID:          testUserID,
	}
}
