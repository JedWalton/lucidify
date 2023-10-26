// //go:build integration
// // +build integration
package syncapi

import (
	"io/ioutil"
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

	// Invalid endpoint
	// Make a request to the /api/sync/localstorage endpoint
	res, err := http.Get(server.URL + "/api/sync/localstorage/?key=test")
	if err != nil {
		t.Fatalf("Failed to make a GET request: %v", err)
	}
	defer res.Body.Close()

	// Reading the response body
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedResponse := `{"success":false,"message":"Invalid key"}` + "\n"
	actualResponse := string(bodyBytes)
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body to be %v; got %v", expectedResponse, string(bodyBytes))
	}

	// Asserting status code and possibly response body (if needed)
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("bodybye: %v", string(bodyBytes))
		t.Fatalf("Expected status OK; got %v", res.StatusCode)
	}

	// Valid endpoint
	// Make a request to the /api/sync/localstorage endpoint
	res, err = http.Get(server.URL + "/api/sync/localstorage/?key=apiKey")
	if err != nil {
		t.Fatalf("Failed to make a GET request: %v", err)
	}
	defer res.Body.Close()

	// Reading the response body
	bodyBytes, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedResponse = `{"success":true,"message":"Successful Get placeholder for key: apiKey"}` + "\n"
	actualResponse = string(bodyBytes)
	if actualResponse != expectedResponse {
		t.Fatalf("Expected response body to be %v; got %v", expectedResponse, string(bodyBytes))
	}

	// Asserting status code and possibly response body (if needed)
	if res.StatusCode != http.StatusOK {
		t.Errorf("bodybye: %v", string(bodyBytes))
		t.Fatalf("Expected status OK; got %v", res.StatusCode)
	}

	// If you want to check the response body as well, you can read and assert it here
	// body, _ := io.ReadAll(res.Body)
	// Perform assertions based on the expected response body
	// e.g., assert.Equal(t, "ExpectedResponse", string(body))
}

func TestHandlerIntegration(t *testing.T) {

}
