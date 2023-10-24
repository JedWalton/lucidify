package chatapi

import (
	"encoding/json"
	"fmt"
	"lucidify-api/service/chatservice"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func ChatHandler(clerkInstance clerk.Client, chatService chatservice.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// How to get currently active user id from clerk
		// ctx := r.Context()
		//
		// sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
		// if !ok {
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	w.Write([]byte("Unauthorized"))
		// 	return
		// }
		//
		// user, err := clerkInstance.Users().Read(sessClaims.Claims.Subject)
		// if err != nil {
		// 	panic(err)
		// }

		// w.Write([]byte(*&user.ID))

		var reqBody map[string]string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		userPrompt := reqBody["message"]
		fmt.Printf("User prompt: %s\n", userPrompt)

		// Do something with the user prompt here
		// CreateWeaviateClass()
		// responseMessage, err := chatService.ProcessCurrentThreadAndReturnSystemPrompt()
		// if err != nil {
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }
		//
		// responseBody := map[string]string{
		// 	"response": responseMessage,
		// }

		w.Header().Set("Content-Type", "application/json")
	}
}

func SyncHandler(chatService chatservice.ChatService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Set content type for all responses from this handler

		switch r.Method {
		case http.MethodGet:
			// Logic for fetching data by key, replace with actual logic
			data, err := fetchDataFromDB("exampleKey") // this function should be implemented
			if err != nil {
				http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
				return
			}

			response := map[string]interface{}{
				"status": "success",
				"data":   data,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)

		case http.MethodPost:
			// Logic for saving/updating data
			SyncData(w, r) // this function should be implemented and return an error if it fails
			// err := SyncData(w, r) // this function should be implemented and return an error if it fails
			// if err != nil {
			// 	http.Error(w, "Failed to sync data", http.StatusInternalServerError)
			// 	return
			// }

			response := map[string]string{
				"status":  "success",
				"message": "Data synced successfully",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)

		case http.MethodDelete:
			// Logic for deleting data by key
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(w, "Key not provided", http.StatusBadRequest)
				return
			}

			DeleteData(w, r, key) // this function should be implemented and return an error if it fails
			// err := DeleteData(w, r, key) // this function should be implemented and return an error if it fails
			// if err != nil {
			// 	http.Error(w, "Failed to delete data", http.StatusInternalServerError)
			// 	return
			// }

			response := map[string]string{
				"status":  "success",
				"message": "Data deleted successfully",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
