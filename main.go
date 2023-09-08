// main.go

package main

import (
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

	allowedOrigins := []string{"http://127.0.0.1:8000"}
	http.Handle("/chat", loggingMiddleware(CORSMiddleware(allowedOrigins)(chatHandler)))
	// http.Handle("/chat", loggingMiddleware(http.HandlerFunc(chatHandler)))
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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		// Continue to the next handler
		next.ServeHTTP(w, r)
	})
}

func CORSMiddleware(allowedOrigins []string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			for _, allowed := range allowedOrigins {
				if origin == allowed || allowed == "*" {
					w.Header().Set("Access-Control-Allow-Origin", allowed)
					w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
					w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
					break
				}
			}

			// If preflight request, respond appropriately
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}
}
