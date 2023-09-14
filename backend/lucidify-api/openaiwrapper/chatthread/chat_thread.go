// chat_completion_controller.go

package chatthread

import (
	"io"
	"lucidify-api/openaiwrapper/chatmodel"
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

	// Here is where the logic to generate pinecone context would go
	response, err := c.Client.SendMessage(userPrompt, "Ask a question about how we can provide value? Follow the conversation thread to uncover value. if you ask a question, only ask one question at a time. Lead them to articulate why I solve their problem. Close the sale.")
	if err != nil {
		return "Error: " + err.Error()
	}

	return response.Choices[0].Message.Content
}
