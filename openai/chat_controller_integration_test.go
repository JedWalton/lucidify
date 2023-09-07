//go:build integration
// +build integration

package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"openai-integrations/utils"
	"os"
	"strings"
	"testing"
	"time"
)

func TestChatControllerIntegration(t *testing.T) {
	// 1. Simulate user input
	input := "Hello, how can you help my business?\nExit\n"
	buffer := bytes.NewBufferString(input)

	if err := utils.LoadDotEnv(); err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}
	controller := NewChatController(apiKey)
	controller.Input = buffer // Set the simulated input

	// 2. Capture the standard output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// 3. Use context to limit runtime for safety
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	controller.Start(ctx)

	// 4. Capture and verify the output
	w.Close()
	os.Stdout = oldStdout
	out, _ := io.ReadAll(r)

	if !strings.Contains(string(out), "Assistant:") {
		t.Errorf("Expected the Assistant's response, but it was not found.")
	}
}
