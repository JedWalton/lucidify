// main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"openai-integrations/openai"
	"openai-integrations/utils"
	"os"
)

var thread *openai.ChatController

func main() {
	if err := utils.LoadDotEnv(); err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	thread = openai.NewChatThread(OPENAI_API_KEY)

	thread.Start(context.Background())

	http.HandleFunc("/chat", chatHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
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
