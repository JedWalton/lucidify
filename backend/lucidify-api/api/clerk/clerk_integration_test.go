// // go:build integration
// // +build integration
package clerk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"lucidify-api/modules/config"
	"lucidify-api/modules/store"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func createUserInClerk(apiKey string) (string, error) {
	url := "https://api.clerk.dev/v1/users"
	payload := strings.NewReader(`{
        "external_id": "test_external_id",
        "first_name": "Test",
        "last_name": "User",
        "email_address": ["test@example.com"],
        "password": "securePassword123"
    }`)

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Authorization", "Bearer "+apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if userID, ok := result["id"].(string); ok {
		return userID, nil
	}
	return "", fmt.Errorf("Failed to create user in Clerk. Response: %s", string(body))
}

func deleteUserInClerk(apiKey string, userID string) error {
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
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("Failed to update user in Clerk. Status code: %d. Response: %s", res.StatusCode, string(body))
	}

	return nil
}

func updateUserInClerk(apiKey string, userID string, firstName string, lastName string) error {
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
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("Failed to update user in Clerk. Response: %s", string(body))
	}

	return nil
}

func checkUserInDB(db *store.Store, userID string, retries int) error {
	for i := 0; i < retries; i++ {
		_, err := db.GetUser(userID)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	return fmt.Errorf("User not found after %d retries", retries)
}

func TestIntegration_usercreatedevent(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	db := testconfig.TestStore

	clerkSecretKey := testconfig.ClerkSecretKey
	userID, err := createUserInClerk(clerkSecretKey)
	if err != nil {
		t.Fatalf("Failed to create user in Clerk: %v", err)
	}
	log.Printf("Created user in Clerk with userID: %s\n", userID)

	err = checkUserInDB(db, userID, 5) // Try 5 times
	if err != nil {
		t.Fatalf("Should have fetched user, user with ID, userID: %v", err)
	}

	err = deleteUserInClerk(clerkSecretKey, userID)
}

func TestIntegration_usercreatedevent_unauthenticated(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	db := testconfig.TestStore

	// clerkSecretKey := testconfig.ClerkSecretKey
	// userID, err := createUserInClerk(clerkSecretKey)
	// if err != nil {
	// 	t.Fatalf("Failed to create user in Clerk: %v", err)
	// }
	// log.Printf("Created user in Clerk with userID: %s\n", userID)

	MakeCurlRequest := func() (int, string, error) {
		cmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", "-X", "POST", "http://localhost:8080/clerk/webhook", "-H", "Content-Type: application/json", "-d", "@test/example_user_created_event.txt")
		out, err := cmd.Output()
		if err != nil {
			return 0, "", err
		}
		statusCode, _ := strconv.Atoi(string(out))
		return statusCode, string(out), nil
	}

	statusCode, response, err := MakeCurlRequest()
	if err != nil {
		t.Fatalf("Error making curl request: %v", err)
	}

	// Check if the status code indicates success (e.g., 200 OK)
	if statusCode >= 200 && statusCode < 300 {
		t.Fatalf("Expected the request to fail, but got a %d status code.", statusCode)
	}

	log.Printf("Response: %s\n", response)

	content, err := ioutil.ReadFile("test/example_user_created_event.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse the content
	var event ClerkEvent
	err = json.Unmarshal(content, &event)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Now you can use the event variable
	UserID := event.Data["id"].(string)

	err = checkUserInDB(db, UserID, 3) // Try 5 times
	if err == nil {
		db.DeleteUser(UserID)
		t.Fatalf("Should have failed to fetch user, unauthenticated user not in db with UserID: %v", err)
	}
	db.DeleteUser(UserID)
}

func TestIntegration_UpdateUser(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	db := testconfig.TestStore

	clerkSecretKey := testconfig.ClerkSecretKey
	userID, err := createUserInClerk(clerkSecretKey)
	if err != nil {
		t.Fatalf("Failed to create user in Clerk: %v", err)
	}
	log.Printf("Created user in Clerk with userID: %s\n", userID)

	err = checkUserInDB(db, userID, 5) // Try 5 times
	if err != nil {
		t.Fatalf("Should have fetched user, user with ID, userID: %v", err)
	}

	// Update the user in Clerk
	err = updateUserInClerk(clerkSecretKey, userID, "UpdatedFirstName", "UpdatedLastName")
	if err != nil {
		t.Fatalf("Failed to update user in Clerk: %v", err)
	}

	time.Sleep(time.Second * 5) // Wait for 5 seconds for the webhook to be processed

	// Check if the user was updated in the local database
	updatedUser, err := db.GetUser(userID)
	if err != nil {
		t.Fatalf("Failed to fetch updated user from local database: %v", err)
	}

	if updatedUser.FirstName != "UpdatedFirstName" || updatedUser.LastName != "UpdatedLastName" {
		t.Fatalf("User update in local database failed. Expected first name: UpdatedFirstName, got: %v. Expected last name: UpdatedLastName, got: %v", updatedUser.FirstName, updatedUser.LastName)
	}

	// Cleanup: Delete the user from Clerk
	err = deleteUserInClerk(clerkSecretKey, userID)
	if err != nil {
		t.Fatalf("Failed to delete user from Clerk: %v", err)
	}
}

func checkUserDeletedFromDB(db *store.Store, userID string, retries int) error {
	for i := 0; i < retries; i++ {
		_, err := db.GetUser(userID)
		if err != nil {
			// If the user is not found, it means the user has been deleted
			return nil
		}
		time.Sleep(time.Second) // Wait for 1 second before retrying
	}
	return fmt.Errorf("User still exists in the database after %d retries", retries)
}

func TestIntegration_userdeletedevent(t *testing.T) {
	testconfig := config.NewTestServerConfig()
	db := testconfig.TestStore

	MakeCurlRequest := func() (string, error) {
		cmd := exec.Command("curl", "-s", "-X", "POST", "http://localhost:8080/clerk/webhook", "-H", "Content-Type: application/json", "-d", "@test/example_user_created_event.txt")
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	response, err := MakeCurlRequest()
	if err != nil {
		t.Fatalf("Failed to make curl request: %v", err)
	}
	log.Printf("Response: %s\n", response)

	content, err := ioutil.ReadFile("test/example_user_created_event.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Parse the content
	var event ClerkEvent
	err = json.Unmarshal(content, &event)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Now you can use the event variable
	UserID := event.Data["id"].(string)

	err = checkUserInDB(db, UserID, 5) // Try 5 times
	if err != nil {
		t.Fatalf("Failed to fetch user, user not in db with UserID: %v", err)
	}

	MakeCurlRequest = func() (string, error) {
		cmd := exec.Command("curl", "-s", "-X", "POST", "http://localhost:8080/clerk/webhook", "-H", "Content-Type: application/json", "-d", "@test/example_user_deleted_event.txt")
		out, err := cmd.Output()
		if err != nil {
			return "", err
		}
		return string(out), nil
	}

	_, err = MakeCurlRequest()

	err = checkUserDeletedFromDB(db, UserID, 5) // Try 5 times

	db.DeleteUser(UserID)
}
