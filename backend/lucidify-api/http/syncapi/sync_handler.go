package syncapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"lucidify-api/service/syncservice"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

// LocalStorageKey defines valid keys for LocalStorage operations.
type LocalStorageKey string

const (
	conversationHistory LocalStorageKey = "conversationHistory"
	folders             LocalStorageKey = "folders"
	prompts             LocalStorageKey = "prompts"
	clearConversations  LocalStorageKey = "clearConversations"
)

// IsValid checks if the provided key is a valid LocalStorageKey.
func (key LocalStorageKey) IsValid() bool {
	switch key {
	case conversationHistory, folders, prompts, clearConversations:
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

func SyncHandler(syncService syncservice.SyncService, clerkInstance clerk.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		key := r.URL.Query().Get("key")
		if !LocalStorageKey(key).IsValid() {
			sendError(w, "Invalid key", http.StatusBadRequest)
			return
		}

		userID, err := getUserIDFromSession(r, clerkInstance)
		if err != nil {
			sendError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			sendError(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var response syncservice.ServerResponse
		value := string(bodyBytes)

		switch r.Method {
		case http.MethodGet:
			response = syncService.HandleGet(userID, key)
		case http.MethodDelete:
			if key == string(clearConversations) {
				response = syncService.HandleClearConversations(userID)
			} else {
				sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
		case http.MethodPost:
			response = syncService.HandleSet(userID, key, value)
		default:
			sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sendJSONResponse(w, http.StatusOK, response)
	})
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	response := syncservice.ServerResponse{
		Success: false,
		Message: message,
	}
	sendJSONResponse(w, statusCode, response)
}

func getUserIDFromSession(r *http.Request, clerkInstance clerk.Client) (string, error) {
	ctx := r.Context()
	sessClaims, ok := clerk.SessionFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("session not found")
	}

	user, err := clerkInstance.Users().Read(sessClaims.Claims.Subject)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}
