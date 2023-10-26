// //go:build integration
// // +build integration
package syncapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"lucidify-api/server/config"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/syncservice"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSyncRoutes(t *testing.T) {
	// Mock configuration and Clerk client for SetupRoutes function
	cfg := &config.ServerConfig{}
	clerkClient, err := clerkservice.NewClerkClient()
	if err != nil {
		t.Fatal(err)
	}

	// Setup routes
	mux := http.NewServeMux()
	mux = SetupRoutes(cfg, mux, clerkClient)

	server := httptest.NewServer(mux)
	defer server.Close()

	t.Run("Test SyncHandler GET route", func(t *testing.T) {
		resp, err := http.Get(server.URL + "/api/sync/localstorage/?key=apiKey")
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var response syncservice.ServerResponse
		if err = json.Unmarshal(body, &response); err != nil {
			t.Fatal(err)
		}

		if !response.Success {
			t.Error("Expected success response")
		}
		// Add any other assertions based on expected response data
	})

	t.Run("Test SyncHandler POST route", func(t *testing.T) {
		data := map[string]interface{}{
			// Mock some data payload for POST request
			//...
		}
		jsonData, _ := json.Marshal(data)
		resp, err := http.Post(server.URL+"/api/sync/localstorage/?key=apiKey", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var response syncservice.ServerResponse
		if err = json.Unmarshal(body, &response); err != nil {
			t.Fatal(err)
		}

		if !response.Success {
			t.Error("Expected success response")
		}
		// Add any other assertions based on expected response data
	})

	t.Run("Test ChangeLogHandler POST route", func(t *testing.T) {
		changelog := []ChangeLog{
			{
				Key:       "apiKey",
				Operation: "update",
				Timestamp: 1635220192,
			},
		}
		jsonData, _ := json.Marshal(changelog)
		resp, err := http.Post(server.URL+"/api/sync/changelog", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		fmt.Println(string(body))
		var response syncservice.ServerResponse
		if err = json.Unmarshal(body, &response); err != nil {
			t.Fatal(err)
		}

		if !response.Success {
			t.Error("Expected success response")
		}

		if response.Message != "Changelog stored successfully" {
			t.Errorf("Expected message 'Changelog stored successfully', got '%s'", response.Message)
		}
	})
	// Add similar test cases for other HTTP methods and routes if necessary
}
