package chatservice

import (
	"context"
	"fmt"
	"lucidify-api/modules/documentservice"
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"

	"github.com/sashabaranov/go-openai"
)

type ChatService interface {
	ProcessCurrentThreadAndReturnSystemPrompt() (string, error)
	ChatCompletion(string) (string, error)
	PerformVectorSearch(string, string) string
	GetAnswerFromFiles(string, string) (string, error)
}

type ChatServiceImpl struct {
	postgresqlDB    postgresqlclient.PostgreSQL
	weaviateDB      weaviateclient.WeaviateClient
	openaiClient    openai.Client
	documentService documentservice.DocumentService
}

func NewChatService(
	postgresqlDB *postgresqlclient.PostgreSQL,
	weaviateDB weaviateclient.WeaviateClient,
	openaiClient *openai.Client,
	documentService documentservice.DocumentService) ChatService {
	return &ChatServiceImpl{postgresqlDB: *postgresqlDB, weaviateDB: weaviateDB, openaiClient: *openaiClient, documentService: documentService}
}

func (c *ChatServiceImpl) ProcessCurrentThreadAndReturnSystemPrompt() (string, error) {
	// UpdateDatabaseWithCurrentChatThread()
	// performVectorDatabaseSearchOnCurrentThread()
	// generateOptimalSystemPrompt()
	return "PLACEHOLDER RESPONSE", nil
}

func (c *ChatServiceImpl) ChatCompletion(userID string) (string, error) {
	// This is a placeholder for the openai completion function
	// This will return the system prompt
	resp, err := c.openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: c.PerformVectorSearch("I want you to talk about cats in every response", userID),
				},
				// {
				// 	Role:    openai.ChatMessageRoleUser,
				// 	Content: "Hello!",
				// },
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *ChatServiceImpl) GetAnswerFromFiles(question string, userID string) (string, error) {
	// Get the vector embedding for the question. You'll need a Go function equivalent to 'get_embedding' in Python.
	concepts := []string{question}
	TOP_K := 2

	// Query your vector database
	results, err := c.weaviateDB.SearchDocumentsByText(TOP_K, userID, concepts)
	if err != nil {
		return "", err
	}

	// Process the results and construct the files_string, similar to Python code
	filesString := ""
	for _, result := range results {
		document, err := c.documentService.GetDocumentByID(userID, result.DocumentID)
		filename := document.DocumentName
		if err != nil {
			return "", err
		}
		fileText := result.ChunkContent
		if result.Certainty > 0.5 {
			fileString := fmt.Sprintf("###\n\"%s\"\n%s\n", filename, fileText)
			filesString += fileString
		}
	}

	// Construct the system message
	systemMessage := fmt.Sprintf(`Given a question, try to answer it using the content of the file extracts below, and if you cannot answer, or find `+
		`a relevant file, just output "I couldn't find the answer to that question in your files.".`+
		`If the answer is not contained in the files or if there are no file extracts, respond with "I couldn't find the answer `+
		`to that question in your files." If the question is not actually a question, respond with "That's not a valid question."`+
		`In the cases where you can find the answer, first give the answer. Then explain how you found the answer from the source or sources, `+
		`and use the exact filenames of the source files you mention. Do not make up the names of any other files other than those mentioned `+
		`in the files context. Give the answer in markdown format.`+
		`Use the following format:`+
		`Question: %s`+
		`Files: %s`+
		`Answer:`, question, filesString)

	// Construct the messages
	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: systemMessage,
		},
		{
			Role:    "user",
			Content: question,
		},
	}

	// Get the completion from OpenAI
	resp, err := c.openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: messages,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// func (c *ChatServiceImpl) PerformVectorSearch(chatCompletionMessage []openai.ChatCompletionMessage, userID string) (string, error) {
func (c *ChatServiceImpl) PerformVectorSearch(chatCompletionMessage string, userID string) string {
	// This is a placeholder for the vector search function
	// This will return the system prompt

	concepts := []string{"Cats"}
	response, err := c.weaviateDB.SearchDocumentsByText(1, userID, concepts)
	if err != nil {
		return "something went fuckin wrong m8 with my apache helicopter"
	}
	return response[0].ChunkContent
	// return "PLACEHOLDER RESPONSE"
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
