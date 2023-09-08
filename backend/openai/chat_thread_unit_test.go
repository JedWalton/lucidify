package openai

import (
	"bufio"
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

func TestProcessUserPrompt(t *testing.T) {
	controller := &ChatController{
		Client:  &FakeClient{},
		Scanner: bufio.NewScanner(strings.NewReader("Test prompt")),
	}

	response := controller.ProcessUserPrompt("Test prompt")
	expectedResponse := "Fake Response"

	if response != expectedResponse {
		t.Fatalf("expected %q, got %q", expectedResponse, response)
	}
}
