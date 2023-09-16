package clerk

import (
	"encoding/json"
	"log"
	"net/http"
)

type ClerkEvent struct {
	Data   map[string]interface{} `json:"data"`
	Object string                 `json:"object"`
	Type   string                 `json:"type"`
}

func ClerkHandler(w http.ResponseWriter, r *http.Request) {
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
		// Handle user created event
		log.Printf("User created: %+v", event.Data)
	case "user.updated":
		// Handle user updated event
		log.Printf("User updated: %+v", event.Data)
	case "user.deleted":
		// Handle user deleted event
		log.Printf("User deleted: %+v", event.Data)
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received"))
}
