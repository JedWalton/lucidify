// //go:build integration
// // +build integration
package syncservice

import "testing"

func TestGetHandlerIntegration(t *testing.T) {
	// Test that the GetHandler returns the expected response
	data, responseObj := HandleGet("apiKey")
	if !responseObj.Success {
		t.Error("Expected success response")
	}
	if data != nil {
		t.Error("Expected nil data")
	}
}

func TestSetHandlerIntegration(t *testing.T) {
	response := HandleSet("apiKey", "test")
	if !response.Success {
		t.Error("Expected success response")
	}
}

func TestRemoveHandlerIntegration(t *testing.T) {
	response := HandleRemove("apiKey")
	if response.Success {
		t.Error("Expected failure response")
	}
}
