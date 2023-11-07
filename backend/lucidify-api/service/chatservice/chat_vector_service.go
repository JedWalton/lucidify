package chatservice

import (
	"fmt"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/service/documentservice"

	"github.com/sashabaranov/go-openai"
)

type ChatVectorService interface {
	ConstructSystemMessage(string, string) (string, error)
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

func (c *ChatVectorServiceImpl) ConstructSystemMessage(question string, userID string) (string, error) {
	// Get the vector embedding for the question. You'll need a Go function equivalent to 'get_embedding' in Python.
	concepts := []string{question}
	TOP_K := 4

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
	return systemMessage, nil
}

// func (c *ChatVectorServiceImpl) ConstructSystemMessage(question string, userID string) (string, error) {
// 	// Safety checks
// 	if c == nil {
// 		return "", fmt.Errorf("ChatVectorServiceImpl is nil")
// 	}
// 	if c.weaviateDB == nil {
// 		return "", fmt.Errorf("WeaviateDB is nil")
// 	}
// 	if c.documentService == nil {
// 		return "", fmt.Errorf("DocumentService is nil")
// 	}
//
// 	// Get the vector embedding for the question.
// 	concepts := []string{question}
// 	TOP_K := 4
//
// 	// Query your vector database
// 	results, err := c.weaviateDB.SearchDocumentsByText(TOP_K, userID, concepts)
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(results) == 0 || results == nil {
// 		return "", fmt.Errorf("SearchDocumentsByText returned nil results")
// 	}
//
// 	// Process the results and construct the files_string
// 	filesString := ""
// 	if len(results) == 0 {
// 		filesString = "I couldn't find the answer to that question in your files."
// 	} else {
// 		for _, result := range results {
// 			document, err := c.documentService.GetDocumentByID(userID, result.DocumentID)
// 			if err != nil {
// 				return "", err
// 			}
// 			if document == nil {
// 				return "", fmt.Errorf("GetDocumentByID returned a nil document for ID: %s", result.DocumentID)
// 			}
//
// 			if result.Certainty > 0.72 {
// 				fileString := fmt.Sprintf("###\n\"%s\"\n%s\n", document.DocumentName, result.ChunkContent)
// 				filesString += fileString
// 			}
// 		}
// 	}
//
// 	// Construct the system message
// 	systemMessage := fmt.Sprintf(`Given a question, try to answer it using the content of the file extracts below, and if you cannot answer, or find `+
// 		`a relevant file, just output "I couldn't find the answer to that question in your files.".`+
// 		`If the answer is not contained in the files or if there are no file extracts, respond with "I couldn't find the answer `+
// 		`to that question in your files." If the question is not actually a question, respond with "That's not a valid question."`+
// 		`In the cases where you can find the answer, first give the answer. Then explain how you found the answer from the source or sources, `+
// 		`and use the exact filenames of the source files you mention. Do not make up the names of any other files other than those mentioned `+
// 		`in the files context. Give the answer in markdown format.`+
// 		`Use the following format:`+
// 		`Question: %s`+
// 		`Files: %s`+
// 		`Answer:`, question, filesString)
//
// 	return systemMessage, nil
// }
