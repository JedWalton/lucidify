// //go:build integration
// // +build integration
package syncapi

import (
	"lucidify-api/server/config"
	"lucidify-api/service/clerkservice"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegrationRoutes(t *testing.T) {
	// Setup the routes and server
	cfg := &config.ServerConfig{} // Assuming default values or mock values for this test
	clerkInstance, err := clerkservice.NewClerkClient()
	if err != nil {
		t.Fatalf("Failed to create Clerk client: %v", err)
	}
	mux := http.NewServeMux()
	server := httptest.NewServer(SetupRoutes(cfg, mux, clerkInstance))
	defer server.Close()

	// Make a request to the /api/sync/localstorage endpoint
	res, err := http.Get(server.URL + "/api/sync/localstorage/test")
	if err != nil {
		t.Fatalf("Failed to make a GET request: %v", err)
	}
	defer res.Body.Close()

	// Asserting status code and possibly response body (if needed)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", res.StatusCode)
	}

	// If you want to check the response body as well, you can read and assert it here
	// body, _ := io.ReadAll(res.Body)
	// Perform assertions based on the expected response body
	// e.g., assert.Equal(t, "ExpectedResponse", string(body))
}
