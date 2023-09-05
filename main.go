package main

import (
	"fmt"
	"openai-integrations/openai"
	"openai-integrations/utils"
	"os"
)

const openAIEndpoint = "https://api.openai.com/v1/chat/completions"

func main() {
	if err := utils.LoadDotEnv(); err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(OPENAI_API_KEY)

	response, err := client.SendMessage("What is the capital of France?", "I want you to be batman and not shut up about it")
	if err != nil {
		// Handle error
	}
	assistantResponse := response.Choices[0].Message.Content
	fmt.Println(assistantResponse)

	systemInput := processAssistantResponseAndGenerateNextSystemInput(assistantResponse, client)

	// To continue the conversation
	response, err = client.SendMessage("And what currency do they use?", systemInput)
	if err != nil {
		// Handle error
	}
}

func processAssistantResponseAndGenerateNextSystemInput(assistantResponse string, client *openai.Client) string {
	for _, message := range client.GetSession().Messages {
		fmt.Println(message.Role, ":", message.Content)
	}

	return "Please use the sales script"
}
