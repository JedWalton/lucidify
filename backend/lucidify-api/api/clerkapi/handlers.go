package clerkapi

import (
	"encoding/json"
	"log"
	"lucidify-api/modules/store/postgresqlclient"
	"net/http"
)

type ClerkEvent struct {
	Data   map[string]interface{} `json:"data"`
	Object string                 `json:"object"`
	Type   string                 `json:"type"`
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if strVal, ok := val.(string); ok {
			return strVal
		}
	}
	return ""
}

func getBoolFromMap(m map[string]interface{}, key string) bool {
	if val, ok := m[key]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

func getInt64FromMap(m map[string]interface{}, key string) int64 {
	if val, ok := m[key]; ok {
		if int64Val, ok := val.(int64); ok {
			return int64Val
		}
	}
	return 0
}

func ClerkHandler(db *postgresqlclient.PostgreSQL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var event ClerkEvent
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&event)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		switch event.Type {
		case "user.created":
			user := postgresqlclient.User{
				UserID:           getStringFromMap(event.Data, "id"),
				ExternalID:       getStringFromMap(event.Data, "external_id"),
				Username:         getStringFromMap(event.Data, "username"),
				PasswordEnabled:  getBoolFromMap(event.Data, "password_enabled"),
				Email:            event.Data["email_addresses"].([]interface{})[0].(map[string]interface{})["email_address"].(string),
				FirstName:        getStringFromMap(event.Data, "first_name"),
				LastName:         getStringFromMap(event.Data, "last_name"),
				ImageURL:         getStringFromMap(event.Data, "image_url"),
				ProfileImageURL:  getStringFromMap(event.Data, "profile_image_url"),
				TwoFactorEnabled: event.Data["two_factor_enabled"].(bool),
				CreatedAt:        getInt64FromMap(event.Data, "created_at"),
				UpdatedAt:        getInt64FromMap(event.Data, "updated_at"),
			}

			// Use userService.CreateUserInUsersTable
			err := db.CreateUserInUsersTable(user)
			if err != nil {
				log.Printf("Error creating user: %v", err)
			}
		case "user.updated":
			user := postgresqlclient.User{
				UserID:           getStringFromMap(event.Data, "id"),
				ExternalID:       getStringFromMap(event.Data, "external_id"),
				Username:         getStringFromMap(event.Data, "username"),
				PasswordEnabled:  getBoolFromMap(event.Data, "password_enabled"),
				Email:            event.Data["email_addresses"].([]interface{})[0].(map[string]interface{})["email_address"].(string),
				FirstName:        getStringFromMap(event.Data, "first_name"),
				LastName:         getStringFromMap(event.Data, "last_name"),
				ImageURL:         getStringFromMap(event.Data, "image_url"),
				ProfileImageURL:  getStringFromMap(event.Data, "profile_image_url"),
				TwoFactorEnabled: event.Data["two_factor_enabled"].(bool),
				CreatedAt:        getInt64FromMap(event.Data, "created_at"),
				UpdatedAt:        getInt64FromMap(event.Data, "updated_at"),
			}

			// Use userService.UpdateUserInUsersTable
			err := db.UpdateUserInUsersTable(user)
			if err != nil {
				log.Printf("Error updating user: %v", err)
			}
		case "user.deleted":
			// Use userService.DeleteUserInUsersTable
			err := db.DeleteUserInUsersTable(event.Data["id"].(string))
			if err != nil {
				log.Printf("Error deleting user: %v", err)
			}
		default:
			log.Printf("Unhandled event type: %s", event.Type)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Received"))
	}
}
