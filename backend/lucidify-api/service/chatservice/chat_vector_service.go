package chatservice

import (
	"context"
	"fmt"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/service/documentservice"

	"github.com/sashabaranov/go-openai"
)

type ChatVectorService interface {
	GetAnswerFromFiles(string, string) (string, error)
}

type ChatVectorServiceImpl struct {
	weaviateDB      weaviateclient.WeaviateClient
	openaiClient    openai.Client
	documentService documentservice.DocumentService
}

func NewChatVectorService(
	weaviateDB weaviateclient.WeaviateClient,
	openaiClient *openai.Client,
	documentService documentservice.DocumentService) ChatVectorService {
	return &ChatVectorServiceImpl{weaviateDB: weaviateDB, openaiClient: *openaiClient, documentService: documentService}
}

func (c *ChatVectorServiceImpl) GetAnswerFromFiles(question string, userID string) (string, error) {
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
		if result.Certainty > 0.72 {
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
