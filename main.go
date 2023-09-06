package main

import (
	"bufio"
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
	// client := openai.NewClient(OPENAI_API_KEY)
	client2 := openai.NewClient(OPENAI_API_KEY)

	// userPrompt := "We are struggling generating leads from our sales content."
	var user2Prompt string

	var i int = 0
	for {
		println("Iteration: ", i)
		println("")

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Please enter your prompt: ")

		// Read input from the terminal
		if scanner.Scan() {
			user2Prompt = scanner.Text()
		} else if err := scanner.Err(); err != nil {
			println("Error reading input: ", err)
			continue
		}

		fmt.Println("")
		// response_client, err := client.SendMessage(userPrompt, "This is the customers input")
		// if err != nil {
		// 	println("Error: ", err)
		// }
		// println("Client 1: \n" + response_client.Choices[0].Message.Content)
		// println("")
		// user2Prompt = response_client.Choices[0].Message.Content
		// response_client2, err := client2.SendMessage(user2Prompt, "Ask a question about how this can provide value? Follow the conversation thread to uncover value. Lead them to articulate why they want to buy my microsaas product. Close the sale.")
		response_client2, err := client2.SendMessage(user2Prompt, "Ask a question about how we can provide value? Follow the conversation thread to uncover value. if you ask a question, only ask one question at a time. Lead them to articulate why I solve their problem. Close the sale.")
		if err != nil {
			println("Error: ", err)
		}
		println("Client 2: \n" + response_client2.Choices[0].Message.Content)
		// userPrompt = response_client2.Choices[0].Message.Content
		println("")

		i++
	}
}
