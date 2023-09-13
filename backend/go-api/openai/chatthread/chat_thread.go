// chat_completion_controller.go

package chatthread

import (
	"encoding/json"
	"io"
	"net/http"
	"openai-integrations/openai/chatmodel"
	"os"
)

type ChatClient interface {
	SendMessage(userInput string, systemInput string) (*chatmodel.CompletionResponse, error)
}

type ChatController struct {
	Client ChatClient
	Input  io.Reader
}

func NewChatThread(apiKey string) *ChatController {
	return &ChatController{
		Client: chatmodel.NewClient(apiKey),
		Input:  os.Stdin,
	}
}

func (c *ChatController) ProcessUserPrompt(userPrompt string) string {

	response, err := c.Client.SendMessage(userPrompt, "Ask a question about how we can provide value? Follow the conversation thread to uncover value. if you ask a question, only ask one question at a time. Lead them to articulate why I solve their problem. Close the sale.")
	if err != nil {
		return "Error: " + err.Error()
	}

	// if isOptimalMomentToScheduleSalesCall() {
	// 	return "Optimal moment to schedule a sales call..."
	// }

	return response.Choices[0].Message.Content
}

func (c *ChatController) chatHandler() http.HandlerFunc {
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
		responseMessage := c.Client.ProcessUserPrompt(userPrompt)

		responseBody := map[string]string{
			"response": responseMessage,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseBody)
	}
}
