// chat_completion_controller.go

package openai

import (
	"bufio"
	"context"
	"fmt"
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

func (c *ChatController) Start(ctx context.Context) {
	i := 0
	c.Scanner = bufio.NewScanner(c.Input)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("Iteration: ", i)
			fmt.Println()

			fmt.Print("Please enter your prompt: ")

			if c.Scanner.Scan() {
				userPrompt := c.Scanner.Text()
				c.processUserPrompt(userPrompt)
			} else if err := c.Scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			i++
		}
	}
}

func (c *ChatController) processUserPrompt(userPrompt string) {
	println(EstimateTokenCount(userPrompt))

	response, err := c.Client.SendMessage(userPrompt, "Ask a question about how we can provide value? Follow the conversation thread to uncover value. if you ask a question, only ask one question at a time. Lead them to articulate why I solve their problem. Close the sale.")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Assistant:\n" + response.Choices[0].Message.Content)

	if isOptimalMomentToScheduleSalesCall() {
		fmt.Println("Optimal moment to schedule a sales call...")
		scheduleSalesCall()
	}
}
