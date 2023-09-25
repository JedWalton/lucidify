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

func getUserIDByEmail(email string, bearerToken string) (string, error) {
	// Construct the URL
	baseURL := "https://api.clerk.dev/v1"
	url := fmt.Sprintf("%s/users?email_address=%s", baseURL, email)

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Set the Authorization header for bearer authentication
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	// Make the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var users []map[string]interface{}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return "", err
	}

	// Check if we got a user and return the ID
	if len(users) > 0 {
		if id, ok := users[0]["id"].(string); ok {
			return id, nil
		}
		return "", fmt.Errorf("id is not a string or not found")
	}

	return "", fmt.Errorf("no user found with email: %s", email)
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
		return fmt.Errorf("Failed to delete user in Clerk. Status code: %d. Response: %s", res.StatusCode, string(body))
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

func RetrieveUser(apiKey string, userID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.clerk.dev/v1/users/%s", userID)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("Failed to retrieve user from Clerk. Response: %s", string(body))
	}

	var user map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
