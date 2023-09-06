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
	client2 := openai.NewClient(OPENAI_API_KEY)

	var userPrompt string

	var i int = 0
	for {
		println("Iteration: ", i)
		println("")

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Please enter your prompt: ")

		if scanner.Scan() {
			userPrompt = scanner.Text()
		} else if err := scanner.Err(); err != nil {
			println("Error reading input: ", err)
			continue
		}
		fmt.Println("")

		response_client, err := client2.SendMessage(userPrompt, "Ask a question about how we can provide value? Follow the conversation thread to uncover value. if you ask a question, only ask one question at a time. Lead them to articulate why I solve their problem. Close the sale.")
		if err != nil {
			println("Error: ", err)
		}
		println("Client 2: \n" + response_client.Choices[0].Message.Content)

		if isOptimalMomentToScheduleSalesCall() {
			println("Optimal moment to schedule a sales call...")
			scheduleSalesCall()
		}
		println("")

		i++
	}
}

func scheduleSalesCall() {
	panic("unimplemented")
}

func isOptimalMomentToScheduleSalesCall() bool {
	return false
}
