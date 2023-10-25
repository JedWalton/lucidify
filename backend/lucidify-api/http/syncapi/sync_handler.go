package syncapi

import (
	"encoding/json"
	"log"
	"lucidify-api/service/syncservice"
	"net/http"
)

// ServerResponse is the structure that defines the standard response from the server.
type ServerResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Holds the actual data, if any
	Message string      `json:"message,omitempty"` // Descriptive message, especially useful in case of errors
}

// LocalStorageKey defines valid keys for LocalStorage operations.
type LocalStorageKey string

const (
	ApiKey               LocalStorageKey = "apiKey"
	ConversationHistory  LocalStorageKey = "conversationHistory"
	SelectedConversation LocalStorageKey = "selectedConversation"
	Theme                LocalStorageKey = "theme"
	Folders              LocalStorageKey = "folders"
	Prompts              LocalStorageKey = "prompts"
	ShowChatbar          LocalStorageKey = "showChatbar"
	ShowPromptbar        LocalStorageKey = "showPromptbar"
	PluginKeys           LocalStorageKey = "pluginKeys"
	Settings             LocalStorageKey = "settings"
)

// IsValid checks if the provided key is a valid LocalStorageKey.
func (key LocalStorageKey) IsValid() bool {
	switch key {
	case ApiKey, ConversationHistory, SelectedConversation, Theme, Folders,
		Prompts, ShowChatbar, ShowPromptbar, PluginKeys, Settings:
		return true
	}
	return false
}

// This is a utility function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, statusCode int, response ServerResponse) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SyncHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Set content type for all responses from this handler

		log.Printf("Request method: %s, URL: %s", r.Method, r.URL.String())

		switch r.Method {
		case http.MethodGet, http.MethodDelete, http.MethodPost:
			// For GET, DELETE, and POST, read 'key' from query parameters
			// key := r.URL.Query().Get("key")
			key := LocalStorageKey(r.URL.Query().Get("key"))
			// Validate the key
			if !key.IsValid() {
				// Handle invalid key: return an error response, etc.
				return
			}
			if key == "" {
				response := ServerResponse{
					Success: false,
					Message: "Key not provided",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(response)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			if r.Method == http.MethodGet {
				syncservice.FetchData(w, r, string(key))
			} else if r.Method == http.MethodDelete {
				syncservice.DeleteData(w, r, string(key))
			} else {
				// For POST, read 'value' from the request body
				var requestData map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&requestData)
				if err != nil {
					http.Error(w, "Bad request data", http.StatusBadRequest)
					return
				}

				value, valueExists := requestData["value"]
				if !valueExists {
					http.Error(w, "Value not provided in request body", http.StatusBadRequest)
					return
				}

				// Now you have the 'key' from the URL and 'value' from the request body and can proceed
				// You might want to modify your 'syncDataToDB' function to accept both 'key' and 'value'
				err = syncservice.SyncData(string(key), value) // Make sure this function accepts both key and value
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// Return a successful response if there were no errors
				response := ServerResponse{
					Success: true,
					Message: "Data synced successfully",
				}
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(response)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
