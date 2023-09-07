package openai

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

type FakeClient struct{}

func (f *FakeClient) SendMessage(userInput string, systemInput string) (*CompletionResponse, error) {
	return &CompletionResponse{
		Choices: []Choice{
			{
				Message: Message{Content: "Fake Response"},
			},
		},
	}, nil
}

// A helper to capture printed output
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestProcessUserPrompt(t *testing.T) {
	controller := &ChatController{
		Client:  &FakeClient{},
		Scanner: bufio.NewScanner(strings.NewReader("Test prompt")),
	}

	// Capturing output
	output := captureOutput(func() {
		controller.processUserPrompt("Test prompt")
	})

	expectedOutput := "Assistant:\nFake Response\n"

	if output != expectedOutput {
		t.Fatalf("expected %q, got %q", expectedOutput, output)
	}
}
