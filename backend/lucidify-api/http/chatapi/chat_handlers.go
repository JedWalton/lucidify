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
