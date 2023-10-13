package chatservice

import (
	"lucidify-api/modules/store/postgresqlclient"
	"lucidify-api/modules/store/weaviateclient"
)

type ChatService interface {
	ProcessCurrentThreadAndReturnSystemPrompt() (string, error)
}

type ChatServiceImpl struct {
	postgresqlDB postgresqlclient.PostgreSQL
	weaviateDB   weaviateclient.WeaviateClient
}

func NewChatService(
	postgresqlDB *postgresqlclient.PostgreSQL,
	weaviateDB weaviateclient.WeaviateClient) ChatService {
	return &ChatServiceImpl{postgresqlDB: *postgresqlDB, weaviateDB: weaviateDB}
}

func (c *ChatServiceImpl) ProcessCurrentThreadAndReturnSystemPrompt() (string, error) {
	// UpdateDatabaseWithCurrentChatThread()
	// performVectorDatabaseSearchOnCurrentThread()
	// generateOptimalSystemPrompt()
	// generateOptimalSystemPromptContext()
	return "PLACEHOLDER RESPONSE", nil
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

// func performVectorDatabaseSearch()
//		performSearchDocumentByText()
//		processResultsOfVectorSearch()

// func constructSystemPromptWithFilteredResults()

// func constructResponseObject() {
//		return New Response
//}
