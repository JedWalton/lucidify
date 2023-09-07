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
	thread := openai.NewChatThread(OPENAI_API_KEY)

	thread.Start(context.Background())
}
