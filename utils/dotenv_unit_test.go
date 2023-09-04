package utils

import (
	"os"
	"testing"
)

func TestLoadDotEnv(t *testing.T) {
	content := []byte(`
# This is a comment
KEY1=VALUE1

# Another comment
KEY2=VALUE2
`)

	tempDir, err := os.MkdirTemp("", "dot-env-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	envFilePath := tempDir + "/.env"
	if err := os.WriteFile(envFilePath, content, 0666); err != nil {
		t.Fatal(err)
	}

	originalWd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalWd)

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() failed: %s", err)
	}

	if val, _ := os.LookupEnv("KEY1"); val != "VALUE1" {
		t.Errorf("Expected VALUE1 for KEY1, got %s", val)
	}
	if val, _ := os.LookupEnv("KEY2"); val != "VALUE2" {
		t.Errorf("Expected VALUE2 for KEY2, got %s", val)
	}
}
