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
		log.Printf("Request method: %s, URL: %s, RemoteAddr: %s", r.Method, r.URL.String(), r.RemoteAddr)

		switch r.Method {
		case http.MethodGet, http.MethodDelete, http.MethodPost:
			key := LocalStorageKey(r.URL.Query().Get("key"))
			if !key.IsValid() {
				sendJSONResponse(w, http.StatusBadRequest, ServerResponse{
					Success: false,
					Message: "Invalid key provided",
				})
				return
			}
			if key == "" {
				sendJSONResponse(w, http.StatusBadRequest, ServerResponse{
					Success: false,
					Message: "Key not provided",
				})
				return
			}

			switch r.Method {
			case http.MethodGet:
				// syncservice.FetchData(w, r, string(key))
				syncservice.FetchData(w, r, string(key))
			case http.MethodDelete:
				syncservice.DeleteData(w, r, string(key))
			case http.MethodPost:
				var requestData map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&requestData)
				if err != nil {
					sendJSONResponse(w, http.StatusBadRequest, ServerResponse{
						Success: false,
						Message: "Bad request data",
					})
					return
				}

				value, valueExists := requestData["value"]
				if !valueExists {
					sendJSONResponse(w, http.StatusBadRequest, ServerResponse{
						Success: false,
						Message: "Value not provided in request body",
					})
					return
				}

				err = syncservice.SyncData(string(key), value)
				if err != nil {
					sendJSONResponse(w, http.StatusInternalServerError, ServerResponse{
						Success: false,
						Message: err.Error(),
					})
					return
				}

				sendJSONResponse(w, http.StatusOK, ServerResponse{
					Success: true,
					Message: "Data synced successfully",
				})
			}

		default:
			sendJSONResponse(w, http.StatusMethodNotAllowed, ServerResponse{
				Success: false,
				Message: "Method not allowed",
			})
		}
	}
}
