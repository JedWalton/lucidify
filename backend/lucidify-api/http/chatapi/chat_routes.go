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
		log.Printf("jsonResponse: error marshalling payload: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("jsonResponse: error writing response: %v", err)
	}
}

func FetchData(w http.ResponseWriter, r *http.Request, key string) {
	log.Printf("FetchData called with key: %s\n", key)

	data, err := fetchDataFromDB(key) // Ensure this function returns the expected data or error
	if err != nil {
		log.Printf("FetchData: error fetching data from DB: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if data == nil {
		log.Printf("FetchData: no data found for key: %s", key)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, Response{"success", data})
}

func DeleteData(w http.ResponseWriter, r *http.Request, key string) {
	log.Printf("DeleteData called with key: %s\n", key)

	err := deleteDataFromDB(key) // Ensure this function returns error if any
	if err != nil {
		log.Printf("DeleteData: error deleting data from DB: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, Response{"success", "data deleted"})
}

func SyncData(w http.ResponseWriter, r *http.Request) {
	log.Println("SyncData called")

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("SyncData: error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	key, ok := requestBody["key"].(string)
	if !ok {
		log.Println("SyncData: 'key' not in request body or not a string")
		http.Error(w, "'key' not in request body or not a string", http.StatusBadRequest)
		return
	}

	value, ok := requestBody["value"]
	if !ok {
		log.Println("SyncData: 'value' not in request body")
		http.Error(w, "'value' not in request body", http.StatusBadRequest)
		return
	}

	err = syncDataToDB(key, value) // Ensure this function returns error if any
	if err != nil {
		log.Printf("SyncData: error syncing data to DB: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, Response{"success", "data synced"})
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
