package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadDotEnv() error {
	file, err := os.Open(".env")
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
