package chatservice

import (
	"context"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"

	"github.com/sashabaranov/go-openai"
)

type ChatService interface {
	ProcessCurrentThreadAndReturnSystemPrompt() (string, error)
	ChatCompletion() (string, error)
}

type ChatServiceImpl struct {
	postgresqlDB postgresqlclient.PostgreSQL
	weaviateDB   weaviateclient.WeaviateClient
	openaiClient openai.Client
}

func NewChatService(
	postgresqlDB *postgresqlclient.PostgreSQL,
	weaviateDB weaviateclient.WeaviateClient,
	openaiClient *openai.Client) ChatService {
	return &ChatServiceImpl{postgresqlDB: *postgresqlDB, weaviateDB: weaviateDB, openaiClient: *openaiClient}
}

func (c *ChatServiceImpl) ProcessCurrentThreadAndReturnSystemPrompt() (string, error) {
	// UpdateDatabaseWithCurrentChatThread()
	// performVectorDatabaseSearchOnCurrentThread()
	// generateOptimalSystemPrompt()
	return "PLACEHOLDER RESPONSE", nil
}

func (c *ChatServiceImpl) ChatCompletion() (string, error) {
	// This is a placeholder for the openai completion function
	// This will return the system prompt
	resp, err := c.openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

//
// func UpdateDatabaseWithCurrentChatThread() {
// Is a new thread?
// 	Yes:
// 		Construct a new chat thread in the database
//  No:
//	 Update the existing chat thread in the database
// }
//	This needs to construct the equivalent json of the export data function in the chatbot-ui.
//	This must be able be imported into the chatbot-ui.
//	This will be used to maintain history of current chat threads.
// }

// type Response struct {
// 	chatThreadID string
// 	systemPrompt string
// 	documentNames []string
// }

// func generateOptimalSystemPrompt() {
//		generateOptimalPrompt
//		generateOptimalContextToAppendToPrompt
//}

// func performVectorDatabaseSearch()
//		performSearchDocumentByText()
//		processResultsOfVectorSearch()

// func constructSystemPromptWithFilteredResults()

// func constructResponseObject() {
//		return New Response
//}
