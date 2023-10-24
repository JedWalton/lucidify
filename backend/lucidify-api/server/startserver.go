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

	"github.com/gorilla/handlers"
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
	cts := chatservice.NewChatThreadService(postgre)
	chatService := chatservice.NewChatService(cts, cvs)

	SetupRoutes(
		config,
		mux,
		postgre,
		clerk,
		weaviate,
		documentService,
		chatService,
	)

	// Set up CORS middleware with your desired options.
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Adjust this to the origins you want to allow.
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Wrap the original mux with the CORS handler.
	corsEnabledMux := corsHandler(mux)

	// Use the CORS-enabled mux in your server.
	log.Printf("Server starting on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, corsEnabledMux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
	// BasicLogging(config, mux)
}
