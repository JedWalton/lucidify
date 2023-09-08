// chat_completion_controller.go

package openai

import (
	"bufio"
	"io"
	"os"
)

type ChatClient interface {
	SendMessage(userInput string, systemInput string) (*CompletionResponse, error)
}

type ChatController struct {
	Client  ChatClient
	Scanner *bufio.Scanner
	Input   io.Reader
}

func NewChatThread(apiKey string) *ChatController {
	return &ChatController{
		Client:  NewClient(apiKey),
		Scanner: bufio.NewScanner(nil),
		Input:   os.Stdin,
	}
}

func (c *ChatController) ProcessUserPrompt(userPrompt string) string {

	response, err := c.Client.SendMessage(userPrompt, "Ask a question about how we can provide value? Follow the conversation thread to uncover value. if you ask a question, only ask one question at a time. Lead them to articulate why I solve their problem. Close the sale.")
	if err != nil {
		return "Error: " + err.Error()
	}

	if isOptimalMomentToScheduleSalesCall() {
		return "Optimal moment to schedule a sales call..."
	}

	return response.Choices[0].Message.Content
}
