//go:build integration
// +build integration

package clerk

import (
	"os/exec"
	"testing"
)

func TestIntegration_chat(t *testing.T) {
	MakeCurlRequest := func() (string, error) {
		cmd := exec.Command("curl", "-s", "-X", "POST", "http://localhost:8080/clerk/webhook", "-H", "Content-Type: application/json", "-d", "@example_user_created_event.txt")
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	// fileContent, err := testutils.ReadFromFile("response.txt")
	// if err != nil {
	// 	t.Fatalf("Failed to read from file: %v", err)
	// }
	//
	response, err := MakeCurlRequest()
	if err != nil {
		t.Fatalf("Failed to make curl request: %v", err)
	}

	// if response != fileContent {
	// 	t.Fatalf("Expected %q but got %q", fileContent, response)
	// }
}
