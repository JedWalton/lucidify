// //go:build integration
// // +build integration
package syncapi

import (
	"io"
	"lucidify-api/server/config"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/syncservice"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupServer(t *testing.T) *httptest.Server {
	cfg := &config.ServerConfig{}
	clerkInstance, err := clerkservice.NewClerkClient()
	if err != nil {
		t.Fatalf("Failed to create Clerk client: %v", err)
	}
	mux := http.NewServeMux()
	syncService, err := syncservice.NewSyncService()
	if err != nil {
		t.Fatalf("Failed to create SyncService: %v", err)
	}
	return httptest.NewServer(SetupRoutes(cfg, mux, clerkInstance, syncService))
}

func makeGetRequest(t *testing.T, server *httptest.Server, endpoint string) (*http.Response, string) {
	res, err := http.Get(server.URL + endpoint)
	if err != nil {
		t.Fatalf("Failed to make a GET request: %v", err)
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	return res, string(bodyBytes)
}

func TestInvalidEndpoint(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	res, responseBody := makeGetRequest(t, server, "/api/sync/localstorage/?key=test")

	expectedResponse := `{"success":false,"message":"Invalid key"}`
	if responseBody != expectedResponse {
		t.Fatalf("Expected response body to be %v; got %v", expectedResponse, responseBody)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status BAD REQUEST; got %v", res.StatusCode)
	}
}

func TestValidEndpoint(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	res, responseBody := makeGetRequest(t, server, "/api/sync/localstorage/?key=apiKey")

	expectedResponse := `{"success":false,"message":"Data not found for key: apiKey"}`
	if responseBody != expectedResponse {
		t.Fatalf("Expected response body to be %v; got %v", expectedResponse, responseBody)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", res.StatusCode)
	}
}

// Further tests can be added in a similar manner
