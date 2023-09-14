//go:build integration
// +build integration

package lucidifychat

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

func MakeCurlRequest() (string, error) {
	cmd := exec.Command("curl", "-s", "-X", "POST", "http://localhost:8080/lucidifychat", "-H", "Content-Type: application/json", "-d", `{"message": "hello"}`)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func TestIntegration_lucidifychat(t *testing.T) {
	fileContent, err := ReadFromFile("lucidifychat.txt")
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
