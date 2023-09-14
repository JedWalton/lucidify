package openaiwrapper

import (
	"encoding/json"
	"net/http"
)

type ChatController interface {
	ProcessUserPrompt(userPrompt string) string
}

func ChatHandler(thread ChatController) http.HandlerFunc {
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
		// Assuming you have a global ChatController instance named 'thread'
		responseMessage := thread.ProcessUserPrompt(userPrompt)

		responseBody := map[string]string{
			"response": responseMessage,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseBody)
	}
}
