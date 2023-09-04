package openai

import (
	"bytes"
	"net/http"
	"testing"
)

// TestRequestConstructor_Construct tests the construction of a request.
func TestRequestConstructor_Construct(t *testing.T) {
	rc := &requestConstructor{APIKey: "testAPIKey"}
	req, err := rc.construct("test prompt")

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
	_, err := rc.construct(string(bytes.Repeat([]byte("a"), 10e6))) // too large payload

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
