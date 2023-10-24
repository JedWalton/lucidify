package chatapi

import (
	"encoding/json"
	"log"
	"lucidify-api/server/config"
	"lucidify-api/server/middleware"
	"lucidify-api/service/chatservice"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(
	config *config.ServerConfig,
	mux *http.ServeMux,
	chatService chatservice.ChatService,
	clerkInstance clerk.Client) *http.ServeMux {

	mux = SetupChatHandler(config, mux, chatService, clerkInstance)

	mux = SetupSyncHandler(config, chatService, mux)

	return mux
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func jsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response) // This is a simplification; in production code, you'd want to handle the error that Write can produce.
}

func FetchData(w http.ResponseWriter, r *http.Request, key string) {
	// No need for mux.Vars(r) as we're now passing the key directly
	log.Printf("FetchData called with key: %s\n", key)

	// Replace with logic to fetch data by key from your database
	data, err := fetchDataFromDB(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, Response{"success", data})
}

func DeleteData(w http.ResponseWriter, r *http.Request, key string) error {
	// No need for mux.Vars(r) as we're now passing the key directly
	log.Printf("DeleteData called with key: %s\n", key)

	// Replace with logic to delete data by key from your database
	err := deleteDataFromDB(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	jsonResponse(w, http.StatusOK, Response{"success", nil})
	return nil
}

// Handler to sync data
func SyncData(w http.ResponseWriter, r *http.Request) error {
	log.Printf("SyncData called with key")

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	key := requestBody["key"].(string)
	value := requestBody["value"]

	// Replace with logic to store data in your database
	err = syncDataToDB(key, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	jsonResponse(w, http.StatusOK, Response{"success", nil})
	return nil
}

// ... other code ...

// Placeholder functions - replace with actual database interactions
func fetchDataFromDB(key string) (interface{}, error) {
	// Your logic to fetch data from the database
	return nil, nil // Placeholder return value, update with your logic
}

func syncDataToDB(key string, value interface{}) error {
	// Your logic to save data to the database
	return nil // Placeholder return value, update with your logic
}

func deleteDataFromDB(key string) error {
	// Your logic to delete data from the database
	return nil // Placeholder return value, update with your logic
}

func SetupChatHandler(
	config *config.ServerConfig,
	mux *http.ServeMux,
	chatService chatservice.ChatService,
	clerkInstance clerk.Client) *http.ServeMux {

	handler := ChatHandler(clerkInstance, chatService)

	injectActiveSession := clerk.WithSession(clerkInstance)

	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.Handle("/chat", injectActiveSession(handler))
	// mux.Handle("/api/sync", handler)

	return mux
}
func SetupSyncHandler(config *config.ServerConfig, chatService chatservice.ChatService, mux *http.ServeMux) *http.ServeMux {

	handler := SyncHandler(chatService)

	handler = middleware.CORSMiddleware(config.AllowedOrigins)(handler)
	handler = middleware.Logging(handler)

	mux.Handle("/api/sync", handler)

	return mux
}
