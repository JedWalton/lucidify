package syncapi

import (
	"encoding/json"
	"io"
	"log"
	"lucidify-api/service/syncservice"
	"net/http"
)

// LocalStorageKey defines valid keys for LocalStorage operations.
type LocalStorageKey string

const (
	apiKey               LocalStorageKey = "apiKey"
	conversationHistory  LocalStorageKey = "conversationHistory"
	selectedConversation LocalStorageKey = "selectedConversation"
	theme                LocalStorageKey = "theme"
	folders              LocalStorageKey = "folders"
	prompts              LocalStorageKey = "prompts"
	showChatbar          LocalStorageKey = "showChatbar"
	showPromptbar        LocalStorageKey = "showPromptbar"
	pluginKeys           LocalStorageKey = "pluginKeys"
	settings             LocalStorageKey = "settings"
)

// IsValid checks if the provided key is a valid LocalStorageKey.
func (key LocalStorageKey) IsValid() bool {
	switch key {
	case apiKey, conversationHistory, selectedConversation, theme, folders,
		prompts, showChatbar, showPromptbar, pluginKeys, settings:
		return true
	}
	return false
}

func MethodNotAllowed(w http.ResponseWriter) {
	response := syncservice.ServerResponse{
		Success: false,
		Message: "Method not allowed",
	}
	sendJSONResponse(w, http.StatusMethodNotAllowed, response)
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, response syncservice.ServerResponse) {
	w.WriteHeader(statusCode)
	responseBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(responseBytes)
	if err != nil {
		log.Println("Error writing response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SyncHandler(syncService syncservice.SyncService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		key := r.URL.Query().Get("key")
		if !LocalStorageKey(key).IsValid() {
			response := syncservice.ServerResponse{
				Success: false,
				Message: "Invalid key",
			}
			sendJSONResponse(w, http.StatusBadRequest, response)
			return
		}

		userID := r.Header.Get("X-User-ID")
		log.Println("Received userID:", userID)

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			// Handle the error, maybe return a response indicating the error.
			return
		}
		value := string(bodyBytes)

		log.Println("Received method:", r.Method)

		var response syncservice.ServerResponse

		switch r.Method {
		case http.MethodGet:
			// response = syncService.HandleGet(userID, key)
			response = syncService.HandleGet(key)
		case http.MethodDelete:
			// response = syncService.HandleRemove(userID, key)
			response = syncService.HandleRemove(key)
		case http.MethodPost:
			// response = syncService.HandleSet(userID, key, value)
			response = syncService.HandleSet(key, value)
		default:
			response = syncservice.ServerResponse{
				Success: false,
				Message: "Method not allowed",
			}
			sendJSONResponse(w, http.StatusMethodNotAllowed, response)
			return
		}
		sendJSONResponse(w, http.StatusOK, response)
	}
}
