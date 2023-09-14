package main

import (
	"log"
	"lucidify-api/middleware"
	"lucidify-api/openaiwrapper/chatthread"
	openaiwrapper "lucidify-api/openaiwrapper/handlers"
	"lucidify-api/store"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type ServerConfig struct {
	OPENAI_API_KEY string
	AllowedOrigins []string
	Port           string
	DBStore        *store.DBStore
	ChatController *chatthread.ChatController
}

func NewServerConfig() *ServerConfig {

	OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
	if OPENAI_API_KEY == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}

	thread := chatthread.NewChatThread(OPENAI_API_KEY)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	allowedOrigins := []string{
		"http://localhost:3000",
	}

	dbStore := store.ConnectToPostgres()

	return &ServerConfig{
		OPENAI_API_KEY: OPENAI_API_KEY,
		AllowedOrigins: allowedOrigins,
		Port:           port,
		DBStore:        dbStore,
		ChatController: thread,
	}
}

func main() {
	config := NewServerConfig()
	defer config.DBStore.Close()

	mux := http.NewServeMux()
	mux = SetupRoutes(config, mux)

	log.Printf("Server starting on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func SetupRoutes(config *ServerConfig, mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("/chat", middleware.Chain(
		openaiwrapper.ChatHandler(config.ChatController),
		middleware.CORSMiddleware(config.AllowedOrigins),
		middleware.Logging,
	))

	return mux
}
