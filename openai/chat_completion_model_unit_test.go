package openai

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "testAPIKey"
	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("Expected non-nil client, but got nil")
	}

	if client.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, but got %s", apiKey, client.APIKey)
	}

	// Type assertion for requestConstructor
	rc, ok := client.requestConstructor.(*requestConstructor)
	if !ok {
		t.Error("Expected requestConstructor of type *requestConstructor")
	} else if rc.APIKey != apiKey {
		t.Errorf("Expected requestConstructor.APIKey to be %s, but got %s", apiKey, rc.APIKey)
	}

	// Type assertion for executor
	ex, ok := client.executor.(*executor)
	if !ok {
		t.Error("Expected executor of type *executor")
	} else if ex.client == nil {
		t.Error("Expected non-nil http.Client in executor, but got nil")
	}
}

func TestClient_SendMessage(t *testing.T) {
	client := &Client{
		requestConstructor: &fakeRequestConstructor{},
		executor:           &fakeExecutor{},
		responseParser:     &fakeResponseParser{},
		session:            ChatSession{},
	}

	resp, err := client.SendMessage("Hello, OpenAI!", "System Command")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.Choices[0].Message.Content != "Hello, World!" {
		t.Errorf("Unexpected response content: %s", resp.Choices[0].Message.Content)
	}

	if len(client.session.Messages) != 3 {
		t.Errorf("Expected 3 messages in session, but got %d", len(client.session.Messages))
	}

	userMessage := client.session.Messages[0]
	if userMessage.Role != "user" || userMessage.Content != "Hello, OpenAI!" {
		t.Errorf("Unexpected user message content or role")
	}

	systemMessage := client.session.Messages[1]
	if systemMessage.Role != "system" || systemMessage.Content != "System Command" {
		t.Errorf("Unexpected system message content or role")
	}

	assistantMessage := client.session.Messages[2]
	if assistantMessage.Role != "assistant" || assistantMessage.Content != "Hello, World!" {
		t.Errorf("Unexpected assistant message content or role")
	}
}

func TestChatSession_AddMessage(t *testing.T) {
	session := &ChatSession{}

	session.AddMessage("user", "Hello")
	if len(session.Messages) != 1 {
		t.Error("Expected 1 message, but got", len(session.Messages))
	}
	if session.Messages[0].Role != "user" {
		t.Errorf("Expected role user, got %s", session.Messages[0].Role)
	}
	if session.Messages[0].Content != "Hello" {
		t.Errorf("Expected content 'Hello', got %s", session.Messages[0].Content)
	}
}

// TestRequestConstructor_Construct tests the construction of a request.
func TestRequestConstructor_Construct(t *testing.T) {
	rc := &requestConstructor{APIKey: "testAPIKey"}

	// Create a slice of chatMessage for the test
	messages := []chatMessage{
		{
			Role:    "user",
			Content: "test prompt",
		},
	}

	req, err := rc.construct(messages) // Now we're passing a slice of chatMessage

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if req.Method != http.MethodPost {
		t.Errorf("Expected method POST, got %s", req.Method)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Unexpected content type: %s", req.Header.Get("Content-Type"))
	}

	if req.Header.Get("Authorization") != "Bearer testAPIKey" {
		t.Errorf("Unexpected authorization: %s", req.Header.Get("Authorization"))
	}
}

// TestResponseParser_Parse tests the response parser's ability to unmarshal the body.
func TestResponseParser_Parse(t *testing.T) {
	rp := &responseParser{}
	body := []byte(`{
		"id": "testID",
		"object": "testObject",
		"created": 1630782143,
		"model": "gpt-3.5-turbo",
		"choices": [{"index": 0, "message": {"role": "system", "content": "Test content"}, "finish_reason": "stop"}],
		"usage": {"prompt_tokens": 10, "completion_tokens": 20, "total_tokens": 30}
	}`)
	resp, err := rp.parse(body)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.ID != "testID" {
		t.Errorf("Expected ID to be testID, got %s", resp.ID)
	}
}

type fakeRequestConstructor struct{}
type fakeExecutor struct{}
type fakeResponseParser struct{}

// Stub the methods of our fake structs
func (f *fakeRequestConstructor) construct(messages []chatMessage) (*http.Request, error) {
	// Create a stubbed http.Request
	req, _ := http.NewRequest(http.MethodGet, "http://fakeurl.com", nil)
	return req, nil
}

func (f *fakeExecutor) execute(req *http.Request) ([]byte, error) {
	return []byte(`{
        "id": "test",
        "object": "completion",
        "created": 1234567890,
        "model": "gpt-4.0-turbo",
        "choices": [{"index": 0, "message": {"role": "assistant", "content": "Hello, World!"}}],
        "usage": {"prompt_tokens": 10, "completion_tokens": 20, "total_tokens": 30}
    }`), nil
}

func (f *fakeResponseParser) parse(respBody []byte) (*CompletionResponse, error) {
	var response CompletionResponse
	err := json.Unmarshal(respBody, &response)
	return &response, err
}

// TestExecutor_Execute checks that the executor handles response errors.
func TestExecutor_Execute(t *testing.T) {
	e := &executor{client: &http.Client{}}

	// This is a dummy endpoint for testing purposes.
	// Since we can't mock, it won't provide the expected response, but we can still verify error handling.
	req, _ := http.NewRequest(http.MethodGet, "https://httpbin.org/status/400", nil)

	_, err := e.execute(req)
	if err == nil {
		t.Errorf("Expected an error for non-200 response, but got nil")
	}
}

func TestConstructError(t *testing.T) {
	rc := &requestConstructor{APIKey: "testAPIKey"}

	// Creating a very large message to generate a too large payload error
	messageContent := string(bytes.Repeat([]byte("a"), 10e6))
	messages := []chatMessage{
		{
			Role:    "user",
			Content: messageContent,
		},
	}

	_, err := rc.construct(messages)

	if err == nil {
		t.Error("Expected error for marshaling large data, but got nil")
	}
}

func TestParseError(t *testing.T) {
	rp := &responseParser{}
	_, err := rp.parse([]byte("not valid json"))

	if err == nil {
		t.Error("Expected error for invalid JSON, but got nil")
	}
}

// We will not be able to test ChatCompletion without mocking. Since it's an integration-heavy function,
// testing it fully requires either an end-to-end test or mocks to simulate OpenAI API responses.
