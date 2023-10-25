package syncapi

import (
	"encoding/json"
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
	apiKey               LocalStorageKey = "apiKey"
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
	case apiKey, ConversationHistory, SelectedConversation, Theme, Folders,
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
		w.Header().Set("Content-Type", "application/json")
		key := r.URL.Query().Get("key")

		// Let syncservice handle all logic, validation, and response generation
		switch r.Method {
		case http.MethodGet:
			syncservice.HandleGet(key)
		case http.MethodDelete:
			syncservice.HandleRemove(key)
		case http.MethodPost:
			syncservice.HandleSet(key, r.Body)
		default:
			syncservice.MethodNotAllowed(w)
		}
	}
}
