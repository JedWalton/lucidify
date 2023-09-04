package openai

import (
	"bytes"
	"net/http"
	"testing"
)

// TestNewClient ensures that the NewClient function initializes all required fields.
func TestNewClient(t *testing.T) {
	apiKey := "testAPIKey"
	client := NewClient(apiKey)

	if client.requestConstructor.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, but got %s", apiKey, client.requestConstructor.APIKey)
	}

	if client.executor.client == nil {
		t.Error("Expected HTTP client to be initialized, but it was nil")
	}

	if client.responseParser == nil {
		t.Error("Expected responseParser to be initialized, but it was nil")
	}
}

// TestConstruct checks if the requestConstructor constructs the request correctly.
func TestConstruct_Client(t *testing.T) {
	prompt := "Hello"
	rc := &requestConstructor{APIKey: "testAPIKey"}
	req, err := rc.construct(prompt)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if req.Method != "POST" {
		t.Errorf("Expected request method to be POST, but got %s", req.Method)
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type header to be set to application/json")
	}

	if req.Header.Get("Authorization") != "Bearer testAPIKey" {
		t.Error("Expected Authorization header to be set to Bearer testAPIKey")
	}

	buffer := new(bytes.Buffer)
	_, _ = buffer.ReadFrom(req.Body)
	bodyStr := buffer.String()
	if bodyStr == "" {
		t.Error("Expected request body to not be empty")
	}

	if !bytes.Contains([]byte(bodyStr), []byte(prompt)) {
		t.Errorf("Expected request body to contain the prompt %s", prompt)
	}
}

// TestExecutor_Execute checks that the executor handles response errors.
// This will actually make a network call to a known endpoint for testing.
func TestExecutor_Execute_Client(t *testing.T) {
	e := &executor{client: &http.Client{}}

	// This is a dummy endpoint for testing purposes.
	// It will provide a 400 response which we can use to test error handling.
	req, _ := http.NewRequest(http.MethodGet, "https://httpbin.org/status/400", nil)

	_, err := e.execute(req)
	if err == nil {
		t.Errorf("Expected an error for non-200 response, but got nil")
	}
}

func TestResponseParser_Parse_Client(t *testing.T) {
	rp := &responseParser{}

	// Construct a valid response payload
	body := []byte(`
{
  "id": "testID",
  "object": "testObject",
  "created": 123456,
  "model": "gpt-3.5-turbo",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "system",
        "content": "Hello!"
      },
      "finish_reason": "stop"
    }
  ],
  "usage": {
    "prompt_tokens": 5,
    "completion_tokens": 10,
    "total_tokens": 15
  }
}
	`)

	response, err := rp.parse(body)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Start validating the parsed response against the manually constructed payload
	if response.ID != "testID" {
		t.Errorf("Expected ID to be 'testID', got '%s'", response.ID)
	}
	if response.Choices[0].Message.Content != "Hello!" {
		t.Errorf("Expected first choice message content to be 'Hello!', got '%s'", response.Choices[0].Message.Content)
	}
}

func TestResponseParser_ParseError(t *testing.T) {
	rp := &responseParser{}

	// Construct an invalid JSON payload
	body := []byte(`{ "invalidJson": }`)

	_, err := rp.parse(body)
	if err == nil {
		t.Fatal("Expected an error due to invalid JSON, but got nil")
	}
}
