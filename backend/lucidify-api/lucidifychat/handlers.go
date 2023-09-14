package lucidifychat

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func LucidifyChatHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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

		responseMessage := "PLACEHOLDER RESPONSE"

		responseBody := map[string]string{
			"response": responseMessage,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseBody)
	}
}
