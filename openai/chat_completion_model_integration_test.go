//go:build integration
// +build integration

package openai

// client_integration_test.go

import (
	"fmt"
	"openai-integrations/utils"
	"os"
	"testing"
	"time"
)

const prompt = "Hello, I'm a human. Are you a human?"
const system = "Hello, you are talking to a human. act nice?"

func TestChatCompletionIntegration(t *testing.T) {
	if err := utils.LoadDotEnv(); err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}

	client := NewClient(apiKey)
	response, err := client.SendMessage(prompt, system)
	if err != nil {
		t.Fatalf("Unexpected error during ChatCompletion: %v", err)
	}

	if response == nil {
		t.Fatal("Expected a response, got nil")
	}

	// Check for ID
	if response.ID == "" {
		t.Fatal("Expected ID to be populated")
	}

	// Check for Model
	if response.Model != "gpt-3.5-turbo-0613" {
		t.Fatalf("Expected model to be 'gpt-3.5-turbo-0613', got %s", response.Model)
	}

	// Check for Choices
	if len(response.Choices) == 0 {
		t.Fatal("Expected at least one choice in the response")
	}

	// Check for the first choice's content
	content := response.Choices[0].Message.Content
	if content == "" {
		t.Fatal("Expected content in the first choice, got empty string")
	}

	// Check for FinishReason in the first choice
	if response.Choices[0].FinishReason == "" {
		t.Fatal("Expected FinishReason in the first choice, got empty string")
	}

	// Check for Role in the first choice's message
	if response.Choices[0].Message.Role == "" {
		t.Fatal("Expected Role in the first choice's message, got empty string")
	}

	// Check for Object
	if response.Object != "chat.completion" {
		t.Fatalf("Expected Object to be 'chat.completion', got %s", response.Object)
	}

	// Check for Created timestamp (assuming it should be recent, e.g., within the last 5 minutes)
	// Note: This is just a rudimentary check and might not be necessary
	if response.Created < (time.Now().Unix() - 300) {
		t.Fatal("Expected a recent Created timestamp")
	}

	// Check for Usage stats
	if response.Usage.PromptTokens == 0 {
		t.Fatal("Expected non-zero PromptTokens")
	}

	if response.Usage.CompletionTokens == 0 {
		t.Fatal("Expected non-zero CompletionTokens")
	}

	if response.Usage.TotalTokens == 0 {
		t.Fatal("Expected non-zero TotalTokens")
	}
}
