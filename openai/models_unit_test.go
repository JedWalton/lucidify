package openai

import (
	"net/http"
	"testing"
)

func TestChatMessageInitialization(t *testing.T) {
	message := chatMessage{
		Role:    "user",
		Content: "Hello, World!",
	}

	if message.Role != "user" {
		t.Errorf("Expected Role 'user', but got %s", message.Role)
	}

	if message.Content != "Hello, World!" {
		t.Errorf("Expected Content 'Hello, World!', but got %s", message.Content)
	}
}

func TestChatCompletionPayloadInitialization(t *testing.T) {
	payload := chatCompletionPayload{
		Model: "gpt-3.5-turbo",
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: "Hello, OpenAI!",
			},
		},
		Temperature: 0.5,
	}

	if payload.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected Model 'gpt-3.5-turbo', but got %s", payload.Model)
	}

	if len(payload.Messages) != 1 || payload.Messages[0].Content != "Hello, OpenAI!" {
		t.Errorf("Expected Messages to have one item with Content 'Hello, OpenAI!', but got %#v", payload.Messages)
	}

	if payload.Temperature != 0.5 {
		t.Errorf("Expected Temperature 0.5, but got %f", payload.Temperature)
	}
}

func TestCompletionResponseInitialization(t *testing.T) {
	completion := CompletionResponse{
		ID:      "1",
		Object:  "response",
		Created: 1630841032,
		Model:   "gpt-3.5-turbo",
		Choices: []Choice{
			{
				Index:        0,
				Message:      Message{Role: "model", Content: "Hi there!"},
				FinishReason: "stop",
			},
		},
		Usage: Usage{PromptTokens: 10, CompletionTokens: 15, TotalTokens: 25},
	}

	if completion.ID != "1" {
		t.Errorf("Expected ID '1', but got %s", completion.ID)
	}

	if len(completion.Choices) != 1 || completion.Choices[0].Message.Content != "Hi there!" {
		t.Errorf("Expected Choices to have one item with Message Content 'Hi there!', but got %#v", completion.Choices)
	}

	if completion.Usage.TotalTokens != 25 {
		t.Errorf("Expected TotalTokens 25, but got %d", completion.Usage.TotalTokens)
	}
}

func TestChoiceInitialization(t *testing.T) {
	choice := Choice{
		Index:        0,
		Message:      Message{Role: "model", Content: "Hello!"},
		FinishReason: "stop",
	}

	if choice.Index != 0 {
		t.Errorf("Expected Index 0, but got %d", choice.Index)
	}

	if choice.Message.Content != "Hello!" {
		t.Errorf("Expected Message Content 'Hello!', but got %s", choice.Message.Content)
	}

	if choice.FinishReason != "stop" {
		t.Errorf("Expected FinishReason 'stop', but got %s", choice.FinishReason)
	}
}

func TestMessageInitialization(t *testing.T) {
	message := Message{
		Role:    "model",
		Content: "Hello!",
	}

	if message.Role != "model" {
		t.Errorf("Expected Role 'model', but got %s", message.Role)
	}

	if message.Content != "Hello!" {
		t.Errorf("Expected Content 'Hello!', but got %s", message.Content)
	}
}

func TestUsageInitialization(t *testing.T) {
	usage := Usage{
		PromptTokens:     10,
		CompletionTokens: 15,
		TotalTokens:      25,
	}

	if usage.PromptTokens != 10 {
		t.Errorf("Expected PromptTokens 10, but got %d", usage.PromptTokens)
	}

	if usage.CompletionTokens != 15 {
		t.Errorf("Expected CompletionTokens 15, but got %d", usage.CompletionTokens)
	}

	if usage.TotalTokens != 25 {
		t.Errorf("Expected TotalTokens 25, but got %d", usage.TotalTokens)
	}
}

func TestRequestConstructorInitialization(t *testing.T) {
	constructor := requestConstructor{
		APIKey: "sample_api_key",
	}

	if constructor.APIKey != "sample_api_key" {
		t.Errorf("Expected APIKey 'sample_api_key', but got %s", constructor.APIKey)
	}
}

func TestResponseParserInitialization(t *testing.T) {
	parser := responseParser{}
	// As this struct is empty, there's no property to assert on.
	// Just checking its existence in this test.
	_ = parser
}

func TestExecutorInitialization(t *testing.T) {
	executor := executor{
		client: &http.Client{},
	}

	if executor.client == nil {
		t.Errorf("Expected client not to be nil")
	}
}
