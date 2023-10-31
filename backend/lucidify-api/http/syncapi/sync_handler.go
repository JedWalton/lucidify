package syncapi

import (
	"encoding/json"
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

//	func SyncHandler(syncService syncservice.SyncService, clerkInstance clerk.Client) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
func SyncHandler(syncService syncservice.SyncService, clerkInstance clerk.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// userID := r.Header.Get("X-User-ID")
		ctx := r.Context()
		log.Println("ctx:", ctx)

		sessClaims, ok := clerk.SessionFromContext(ctx)
		log.Println("sessClaims:", sessClaims)
		log.Println("ok:", ok)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		user, err := clerkInstance.Users().Read(sessClaims.Claims.Subject)
		if err != nil {
			panic(err)
		}

		userID := user.ID

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
			response = syncService.HandleGet(userID, key)
		case http.MethodDelete:
			// response = syncService.HandleRemove(userID, key)
			response = syncService.HandleClearConversations(userID)
		case http.MethodPost:
			response = syncService.HandleSet(userID, key, value)
			// response = syncService.HandleSet(key, value)
		default:
			response = syncservice.ServerResponse{
				Success: false,
				Message: "Method not allowed",
			}
			sendJSONResponse(w, http.StatusMethodNotAllowed, response)
			return
		}
		sendJSONResponse(w, http.StatusOK, response)
	})
}
