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
	client2 := openai.NewClient(OPENAI_API_KEY)

	var userPrompt string
	var user2Prompt string

	var i int = 0
	for {
		println("Iteration: ", i)
		println("")
		response_client, err := client.SendMessage(userPrompt, "Talk about microsaas opportunities with large language models. Ideally keep responses to a just a couple of lines. discuss propositions and constructing offers")
		if err != nil {
			println("Error: ", err)
		}
		println("Client 1: \n" + response_client.Choices[0].Message.Content)
		println("")
		user2Prompt = response_client.Choices[0].Message.Content
		response_client2, err := client2.SendMessage(user2Prompt, "Ask a question about how this can provide value? Follow the conversation thread to uncover value. Lead them to articulate why they want to buy my microsaas product. Close the sale.")
		if err != nil {
			println("Error: ", err)
		}
		println("Client 2: \n" + response_client2.Choices[0].Message.Content)
		userPrompt = response_client2.Choices[0].Message.Content
		println("")

		i++
	}
}
