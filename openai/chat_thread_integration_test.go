//go:build integration
// +build integration

package openai

import (
	"context"
	"fmt"
	"openai-integrations/utils"
	"os"
	"testing"
	"time"
)

func TestChatControllerIntegration(t *testing.T) {
	if err := utils.LoadDotEnv(); err != nil {
		fmt.Println("Error loading .env:", err)
		return
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
