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

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type TestSetup struct {
	Config        *config.ServerConfig
	ClerkInstance clerk.Client
	SyncService   syncservice.SyncService
	Server        *httptest.Server
}

func setupTestServer(t *testing.T) *TestSetup {
	cfg := &config.ServerConfig{}

	clerkInstance, err := clerkservice.NewClerkClient()
	if err != nil {
		t.Fatalf("Failed to create Clerk client: %v", err)
	}

	syncService, err := syncservice.NewSyncService()
	if err != nil {
		t.Fatalf("Failed to create SyncService: %v", err)
	}

	mux := http.NewServeMux()
	server := httptest.NewServer(SetupRoutes(cfg, mux, clerkInstance, syncService))

	return &TestSetup{
		Config:        cfg,
		ClerkInstance: clerkInstance,
		SyncService:   syncService,
		Server:        server,
	}
}

func makeGetRequestNoJWT(t *testing.T, server *httptest.Server, endpoint string) (*http.Response, string) {
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

func TestGetRequestNoJWT(t *testing.T) {
	setup := setupTestServer(t)
	defer setup.Server.Close()

	res, responseBody := makeGetRequestNoJWT(t, setup.Server, "/api/sync/localstorage/?key=test")

	expectedResponse := `couldn't find cookie __session`
	if responseBody != expectedResponse {
		t.Fatalf("Expected response body to be %v; got %v", expectedResponse, responseBody)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("Expected status BAD REQUEST; got %v", res.StatusCode)
	}
}
