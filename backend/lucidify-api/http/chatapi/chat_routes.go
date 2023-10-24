package chatapi

import (
	"encoding/json"
	"fmt"
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

func DeleteData(w http.ResponseWriter, r *http.Request, key string) error {
	log.Printf("DeleteData called with key: %s\n", key)

	err := deleteDataFromDB(key) // Ensure this function returns error if any
	if err != nil {
		log.Printf("DeleteData: error deleting data from DB: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	jsonResponse(w, http.StatusOK, Response{"success", "data deleted"})
	return nil
}

func SyncData(w http.ResponseWriter, r *http.Request) error {
	log.Println("SyncData called")

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Printf("SyncData: error decoding request body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	key, ok := requestBody["key"].(string)
	if !ok {
		log.Println("SyncData: 'key' not in request body or not a string")
		http.Error(w, "'key' not in request body or not a string", http.StatusBadRequest)
		return err
	}

	value, ok := requestBody["value"]
	if !ok {
		log.Println("SyncData: 'value' not in request body")
		http.Error(w, "'value' not in request body", http.StatusBadRequest)
		return err
	}

	err = syncDataToDB(key, value) // Ensure this function returns error if any
	if err != nil {
		log.Printf("SyncData: error syncing data to DB: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	jsonResponse(w, http.StatusOK, Response{"success", "data synced"})
	return nil
}

// ... other code ...

// fetchDataFromDB is a stub function to simulate database fetching.
func fetchDataFromDB(key string) (interface{}, error) {
	// Instead of fetching data from a database, we return a hardcoded value.
	// You should replace this with actual database interaction logic.
	data := map[string]string{
		"exampleKey": "exampleValue",
	}

	if value, exists := data[key]; exists {
		return value, nil
	}

	return nil, fmt.Errorf("no data found for key: %s", key)
}

// deleteDataFromDB is a stub function to simulate database deletion.
func deleteDataFromDB(key string) error {
	// Instead of deleting data from a database, we just log the action and pretend it succeeded.
	// You should replace this with actual database interaction logic.
	log.Printf("Data with key '%s' is supposed to be deleted here.", key)
	return nil // no error, means it was "successful"
}

// syncDataToDB is a stub function to simulate database sync/insert/update.
func syncDataToDB(key string, value interface{}) error {
	// Instead of syncing data to a database, we just log the action and pretend it succeeded.
	// You should replace this with actual database interaction logic.
	log.Printf("Data with key '%s' and value '%v' is supposed to be synced here.", key, value)
	return nil // no error, means it was "successful"
}

func SetupChatHandler(
	config *config.ServerConfig,
	mux *http.ServeMux,
	chatService chatservice.ChatService,
	clerkInstance clerk.Client) *http.ServeMux {

	handler := ChatHandler(clerkInstance, chatService)

	injectActiveSession := clerk.WithSession(clerkInstance)

	handler = middleware.Logging(handler)

	mux.Handle("/chat", injectActiveSession(handler))
	// mux.Handle("/api/sync", handler)

	return mux
}

func SetupSyncHandler(config *config.ServerConfig, chatService chatservice.ChatService, mux *http.ServeMux) *http.ServeMux {

	handler := SyncHandler(chatService)

	handler = middleware.Logging(handler)

	// mux.Handle("/api/sync", handler)
	mux.Handle("/api/sync/", http.StripPrefix("/api/sync/", handler))

	return mux
}
