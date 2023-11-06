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

type ChatResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Use this field to include any data in the case of success
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
		// systemPromptFromVecSearch, err := cvs.ConstructSystemMessage(reqBody.Messages[len(reqBody.Messages)-1].Content, user.ID)
		// if err != nil {
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// 	return
		// }
		//
		// // Placeholder response
		// placeholderResponse := map[string]interface{}{
		// 	"status":       "success",
		// 	"systemPrompt": systemPromptFromVecSearch,
		// }
		//
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusOK) // Set the status code to 200 OK
		// if err := json.NewEncoder(w).Encode(placeholderResponse); err != nil {
		// 	// If there is an error when encoding the response, log it and send a server error status
		// 	fmt.Println("Error encoding placeholder response:", err)
		// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
		// }
		// Create a response object
		response := ChatResponse{}

		systemPromptFromVecSearch, err := cvs.ConstructSystemMessage(reqBody.Messages[len(reqBody.Messages)-1].Content, user.ID)
		if err != nil {
			// Handle the failure by setting the response fields accordingly
			response.Status = "fail"
			response.Message = "Internal server error"
			// You can optionally include more error information in response.Data
			// response.Data = err.Error() // Uncomment if you want to send the error message back to the client

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError) // Set the appropriate status code for an error
			json.NewEncoder(w).Encode(response)           // No need to check error here, we're already in an error state
			return
		}

		// Handle the success case by setting the response fields accordingly
		response.Status = "success"
		response.Message = "System prompt constructed successfully"
		response.Data = systemPromptFromVecSearch

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // Set the status code to 200 OK
		if err := json.NewEncoder(w).Encode(response); err != nil {
			// If there is an error when encoding the response, log it and send a server error status
			fmt.Println("Error encoding response:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}
