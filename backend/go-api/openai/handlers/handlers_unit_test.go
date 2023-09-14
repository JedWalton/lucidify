package openai

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StubChatController struct {
	response string
}

func (s StubChatController) ProcessUserPrompt(userPrompt string) string {
	return s.response
}

func TestChatHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(`{"message": "test"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := ChatHandler(StubChatController{response: "stub response"})

	handler.ServeHTTP(rr, req)

	expected := `{"response":"stub response"}`
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
