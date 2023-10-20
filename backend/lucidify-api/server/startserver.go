package server

import (
	"log"
	"lucidify-api/data/store/postgresqlclient"
	"lucidify-api/data/store/weaviateclient"
	"lucidify-api/server/config"
	"lucidify-api/service/chatservice"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
	"net/http"

	"github.com/sashabaranov/go-openai"
)

func StartServer() {
	config := config.NewServerConfig()

	mux := http.NewServeMux()

	postgre, err := postgresqlclient.NewPostgreSQL()
	if err != nil {
		log.Fatal(err)
	}

	clerk, err := clerkservice.NewClerkClient()
	if err != nil {
		log.Fatal(err)
	}

	weaviate, err := weaviateclient.NewWeaviateClient()
	if err != nil {
		log.Fatal(err)
	}

	documentService := documentservice.NewDocumentService(postgre, weaviate)

	openaiClient := openai.NewClient(config.OPENAI_API_KEY)

	cvs := chatservice.NewChatVectorService(weaviate, openaiClient, documentService)
	chs := chatservice.NewChatHistoryService(postgre)
	chatService := chatservice.NewChatService(chs, cvs)

	SetupRoutes(
		config,
		mux,
		postgre,
		clerk,
		weaviate,
		documentService,
		chatService,
	)

	BasicLogging(config, mux)
}
