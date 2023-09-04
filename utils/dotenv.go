package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var getProjectRoot = func() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFilePath), "../")
}

func LoadDotEnv() error {
	projectRoot := getProjectRoot()

	// Create the path to the .env file in the root directory.
	envPath := filepath.Join(projectRoot, ".env")

	file, err := os.Open(envPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		if strings.HasPrefix(strings.TrimSpace(line), "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid line at %d: %s\n", lineNumber, line)
			continue
		}

		key := parts[0]
		value := parts[1]

		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while scanning .env at line %d: %w", lineNumber, err)
	}

	return nil
}
