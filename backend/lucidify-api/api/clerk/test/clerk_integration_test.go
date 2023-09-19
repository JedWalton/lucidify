//go:build integration
// +build integration

package clerk

import (
	"log"
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

	// file, err := testutils.ReadFromFile("example_user_created_event.txt")
	// if err != nil {
	// 	t.Fatalf("Failed to read from file: %v", err)
	// }
	// log.Printf("File: %s\n", file)

	response, err := MakeCurlRequest()
	if err != nil {
		t.Fatalf("Failed to make curl request: %v", err)
	}
	log.Printf("Response: %s\n", response)

	// if response != fileContent {
	// 	t.Fatalf("Expected %q but got %q", fileContent, response)
	// }
}
