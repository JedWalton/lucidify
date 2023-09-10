//go:build integration
// +build integration

package chatthread

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestChatControllerIntegration(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("No .env file found")
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}
	thread := NewChatThread(apiKey)

	// Use context to limit runtime for safety
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	response := thread.ProcessUserPrompt("Hello, how can you help my business?")

	if response == "" {
		t.Errorf("Expected the Assistant's response, but it was empty.")
	}

}
