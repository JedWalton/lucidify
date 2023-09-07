// main.go

package main

import (
	"context"
	"fmt"
	"openai-integrations/openai"
	"openai-integrations/utils"
	"os"
)

func main() {
	if err := utils.LoadDotEnv(); err != nil {
		fmt.Println("Error loading .env:", err)
		return
	}

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	controller := openai.NewChatController(OPENAI_API_KEY)

	controller.Start(context.Background())
}
