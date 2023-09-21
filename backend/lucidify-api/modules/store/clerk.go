package store

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CreateUserInClerk(apiKey string, firstName string, lastName string, email string, password string) (string, error) {
	url := "https://api.clerk.dev/v1/users"
	payload := strings.NewReader(fmt.Sprintf(`{
        "first_name": "%s",
        "last_name": "%s",
        "email_address": ["%s"],
        "password": "%s"
    }`, firstName, lastName, email, password))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if userID, ok := result["id"].(string); ok {
		return userID, nil
	}
	return "", fmt.Errorf("Failed to create user in Clerk. Response: %s", string(body))
}

func DeleteUserInClerk(apiKey string, userID string) error {
	url := fmt.Sprintf("https://api.clerk.dev/v1/users/%s", userID)

	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Failed to update user in Clerk. Status code: %d. Response: %s", res.StatusCode, string(body))
	}

	return nil
}

func UpdateUserInClerk(apiKey string, userID string, firstName string, lastName string) error {
	url := fmt.Sprintf("https://api.clerk.dev/v1/users/%s", userID)
	payload := strings.NewReader(fmt.Sprintf(`{
        "first_name": "%s",
        "last_name": "%s"
    }`, firstName, lastName))

	req, err := http.NewRequest("PATCH", url, payload)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Failed to update user in Clerk. Response: %s", string(body))
	}

	return nil
}
