//go:build integration
// +build integration

package chat

import (
	"io/ioutil"
	"os/exec"
	"testing"
)

func WriteToFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func ReadFromFile(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func TestIntegration_chat(t *testing.T) {
	MakeCurlRequest := func() (string, error) {
		cmd := exec.Command("curl", "-s", "-X", "POST", "http://localhost:8080/chat", "-H", "Content-Type: application/json", "-d", `{"message": "hello"}`)
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	fileContent, err := ReadFromFile("chat.txt")
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
