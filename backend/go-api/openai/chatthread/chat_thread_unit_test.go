package chatthread

import (
	"openai-integrations/openai/chatmodel"
	"testing"
)

type FakeClient struct{}

func (f *FakeClient) SendMessage(userInput string, systemInput string) (*chatmodel.CompletionResponse, error) {
	return &chatmodel.CompletionResponse{
		Choices: []chatmodel.Choice{
			{
				Message: chatmodel.Message{Content: "Fake Response"},
			},
		},
	}, nil
}

func TestProcessUserPrompt(t *testing.T) {
	controller := &ChatController{
		Client: &FakeClient{},
	}

	response := controller.ProcessUserPrompt("Test prompt")
	expectedResponse := "Fake Response"

	if response != expectedResponse {
		t.Fatalf("expected %q, got %q", expectedResponse, response)
	}
}
