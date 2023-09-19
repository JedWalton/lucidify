//go:build integration
// +build integration

package documents

import (
	"lucidify-api/modules/testutils"
	"os/exec"
	"testing"
)

func TestIntegration_documentsupload(t *testing.T) {
	MakeCurlRequest := func() (string, error) {
		cmd := exec.Command("curl", "-s", "-X", "POST", "http://localhost:8080/documents/upload", "-H", "Content-Type: application/json", "-d", `{"title": "hello", "content": "world"}`)
		// cmd := exec.Command("curl", "-s", "-X", "POST", "http://lucidify.serveo.net/documents/upload", "-H", "Content-Type: application/json", "-d", `{"title": "hello", "content": "world"}`)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	fileContent, err := testutils.ReadFromFile("uploaddocument.txt")
	if err != nil {
		t.Fatalf("Failed to read from file: %v", err)
	}

	response, err := MakeCurlRequest()
	if err != nil {
		t.Fatalf("Failed to make curl request: %v", err)
	}

	if response != fileContent {
		t.Fatalf("Expected %q but got %q", fileContent, response)
	}
}
