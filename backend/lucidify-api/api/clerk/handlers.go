package clerk

import (
	"encoding/json"
	"log"
	"lucidify-api/modules/store"
	"net/http"
)

type ClerkEvent struct {
	Data   map[string]interface{} `json:"data"`
	Object string                 `json:"object"`
	Type   string                 `json:"type"`
}

func ClerkHandler(db *store.Store) http.HandlerFunc {
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
			user := store.User{
				UserID:           event.Data["id"].(string),
				ExternalID:       event.Data["external_id"].(string),
				Username:         event.Data["username"].(string),
				PasswordEnabled:  event.Data["password_enabled"].(bool),
				Email:            event.Data["email_addresses"].([]interface{})[0].(map[string]interface{})["email_address"].(string),
				FirstName:        event.Data["first_name"].(string),
				LastName:         event.Data["last_name"].(string),
				ImageURL:         event.Data["image_url"].(string),
				ProfileImageURL:  event.Data["profile_image_url"].(string),
				TwoFactorEnabled: event.Data["two_factor_enabled"].(bool),
				CreatedAt:        event.Data["created_at"].(int64),
				UpdatedAt:        event.Data["updated_at"].(int64),
				Deleted:          false,
			}

			err := db.CreateUser(user)
			if err != nil {
				log.Printf("Error creating user: %v", err)
			}
		case "user.updated":
			user := store.User{
				UserID:           event.Data["id"].(string),
				ExternalID:       event.Data["external_id"].(string),
				Username:         event.Data["username"].(string),
				PasswordEnabled:  event.Data["password_enabled"].(bool),
				Email:            event.Data["email_addresses"].([]interface{})[0].(map[string]interface{})["email_address"].(string),
				FirstName:        event.Data["first_name"].(string),
				LastName:         event.Data["last_name"].(string),
				ImageURL:         event.Data["image_url"].(string),
				ProfileImageURL:  event.Data["profile_image_url"].(string),
				TwoFactorEnabled: event.Data["two_factor_enabled"].(bool),
				UpdatedAt:        event.Data["updated_at"].(int64),
			}

			err := db.UpdateUser(user)
			if err != nil {
				log.Printf("Error updating user: %v", err)
			}
		case "user.deleted":
			err := db.SetUserDeleted(event.Data["id"].(string))
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
