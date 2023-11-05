package chatapi

import (
	"encoding/json"
	"fmt"
	"log"
	"lucidify-api/service/chatservice"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

type ServerResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Holds the actual data, if any
	Message string      `json:"message,omitempty"` // Descriptive message, especially useful in case of errors
}

type Role string

const (
	RoleUser   Role = "user"
	RoleSystem Role = "system"
	// Add other roles as needed
)

// Message corresponds to the TypeScript interface with a role and content.
type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func ChatHandler(clerkInstance clerk.Client, cvs chatservice.ChatVectorService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// How to get currently active user id from clerk
		ctx := r.Context()

		sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		user, err := clerkInstance.Users().Read(sessClaims.Claims.Subject)
		if err != nil {
			panic(err)
		}

		var reqBody struct {
			Messages []Message `json:"messages"`
		}

		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		log.Printf("User prompt: %s\n", reqBody.Messages)
		systemPromptFromVecSearch, err := cvs.ConstructSystemMessage(reqBody.Messages[len(reqBody.Messages)-1].Content, user.ID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Placeholder response
		placeholderResponse := map[string]interface{}{
			"status":       "success",
			"systemPrompt": systemPromptFromVecSearch,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Set the status code to 200 OK
		if err := json.NewEncoder(w).Encode(placeholderResponse); err != nil {
			// If there is an error when encoding the response, log it and send a server error status
			fmt.Println("Error encoding placeholder response:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
