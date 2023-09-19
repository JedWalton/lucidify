package testutils

import (
	"os"
)

func WriteToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func ReadFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
