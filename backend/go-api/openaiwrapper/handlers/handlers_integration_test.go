//go:build integration
// +build integration

package openaiwrapper

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"openaiwrapper-integrations/openaiwrapper/chatthread"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestChatHandlerIntegration(t *testing.T) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		log.Println("No .env file found")
	}

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		t.Fatal("OPENAI_API_KEY not set, skipping integration test")
	}

	chatController := chatthread.NewChatThread(OPENAI_API_KEY)
	handler := ChatHandler(chatController)

	// Create a request
	req, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(`{"message": "test"}`)))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call ServeHTTP directly
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is in the correct format
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if respContent, ok := response["response"]; !ok {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	} else if respContent == "" {
		t.Errorf("response content is empty")
	} else if respContent == "Error: " {
		t.Errorf("handler returned an error: %v", respContent)
	}
}
