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
	openaiClient := openai.NewClient(OPENAI_API_KEY)

	response, err := openaiClient.ChatCompletion("Hello, I'm a human. Are you a human?")
	if err != nil {
		fmt.Println("Error calling OpenAI:", err)
		return
	}

	fmt.Println(response.Choices[0].Message.Content)
}
