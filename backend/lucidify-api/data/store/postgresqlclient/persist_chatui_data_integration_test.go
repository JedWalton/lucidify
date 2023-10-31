// //go:build integration
// // +build integration
package postgresqlclient

import (
	"lucidify-api/data/store/storemodels"
	"testing"
)

func TestPostgreSQLClientFunctions(t *testing.T) {
	store, err := NewPostgreSQL()
	if err != nil {
		t.Errorf("Failed to create test postgresqlclient: %v", err)
	}

	user := storemodels.User{
		UserID:           "pgclient_integration_test_user_id",
		ExternalID:       "TestChatUIExternalID",
		Username:         "TestChatUIUsername",
		PasswordEnabled:  true,
		Email:            "TestChatUI@example.com",
		FirstName:        "TestChatUICreateTest",
		LastName:         "TestChatUIUser",
		ImageURL:         "https://TestChatUIimageurl.jpg",
		ProfileImageURL:  "https://TestChatUIProfile.jpg",
		TwoFactorEnabled: false,
		CreatedAt:        1654012591514,
		UpdatedAt:        1654012591514,
	}

	err = store.CreateUserInUsersTable(user)
	if err != nil {
		t.Errorf("Failed to create user: %v", err)
	}

	userID := "pgclient_integration_test_user_id"

	t.Run("test set and get data", func(t *testing.T) {
		tests := []struct {
			key          string
			value        string
			expectedData string
			tableName    string
			shouldError  bool
		}{
			{"conversationHistory_testKey", "conversationHistory_testData", "conversationHistory_testData", "conversation_history", false},
			{"folders_testKey", "folders_testData", "folders_testData", "folders", false},
			{"prompts_testKey", "prompts_testData", "prompts_testData", "prompts", false},
			{"invalid_testKey", "invalid_testData", "", "", true},
		}

		for _, test := range tests {
			err := store.SetData(userID, test.key, test.value)
			if (err != nil) != test.shouldError {
				t.Errorf("Expected error: %v, got: %v for key: %s", test.shouldError, err, test.key)
			}

			if !test.shouldError {
				data, err := store.GetData(userID, test.key)
				if err != nil {
					t.Errorf("Error fetching data for key: %s, err: %v", test.key, err)
				}
				if data != test.expectedData {
					t.Errorf("Expected data: %s, got: %s for key: %s", test.expectedData, data, test.key)
				}
			}
		}
	})

	t.Run("test clear conversations", func(t *testing.T) {
		// Pre-setup: Insert some data to ensure clear functionality works.
		err := store.SetData(userID, "conversationHistory_preClear", "testData")
		if err != nil {
			t.Errorf("Pre-setup failed for clear conversations test. Err: %v", err)
			return
		}
		err = store.SetData(userID, "folders_preClear", "testData")
		if err != nil {
			t.Errorf("Pre-setup failed for clear conversations test. Err: %v", err)
			return
		}

		err = store.ClearConversations(userID)
		if err != nil {
			t.Errorf("Failed to clear conversations. Err: %v", err)
			return
		}

		// Check if data is actually cleared.
		data, err := store.GetData(userID, "conversationHistory_preClear")
		if err == nil || data != "" {
			t.Errorf("Data not cleared for key: conversationHistory_preClear")
		}
		data, err = store.GetData(userID, "folders_preClear")
		if err == nil || data != "" {
			t.Errorf("Data not cleared for key: folders_preClear")
		}
	})

	t.Cleanup(func() {
		// Delete the test user
		err = store.DeleteUserInUsersTable("pgclient_integration_test_user_id")
		if err != nil {
			t.Errorf("Failed to delete test user: %v", err)
		}
		// Cleanup inserted test data.
		err := store.ClearConversations(userID)
		if err != nil {
			t.Errorf("Cleanup failed. Unable to clear conversations. Err: %v", err)
		}
	})
}
