package syncservice

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// ServerResponse is the structure that defines the standard response from the server.
type ServerResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Holds the actual data, if any
	Message string      `json:"message,omitempty"` // Descriptive message, especially useful in case of errors
}

// This is a utility function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, statusCode int, response ServerResponse) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func FetchData(w http.ResponseWriter, r *http.Request, key string) {
	log.Printf("FetchData called with key: %s\n", key)

	fetchedData, exists := GetDataFromServerDB(key)
	if !exists {
		response := ServerResponse{
			Success: false,
			Message: "No data found for key: " + key,
		}
		sendJSONResponse(w, http.StatusInternalServerError, response)
		return
	}
	serverResponse := ServerResponse{
		Success: true,
		Data:    fetchedData,
		Message: "Data fetched successfully",
	}

	sendJSONResponse(w, http.StatusOK, serverResponse)
}

func DeleteData(w http.ResponseWriter, r *http.Request, key string) {
	response := ServerResponse{
		Success: true,
		Message: "Data deleted successfully",
	}

	sendJSONResponse(w, http.StatusOK, response)
}

// SyncDataToDb function that accepts a key and value and syncs this data with a database.
// You should replace the contents of this function with actual database interaction logic.
func SyncData(key string, value interface{}) error {
	// This function is currently a stub and does not actually interact with a database.
	// Here, you would write your logic to sync data to your database.
	// This might include SQL statements, or calls to another service, etc.

	SetDataInServerDB(key, value)
	// For demonstration, we'll just print the key and value.
	log.Printf("Syncing data to DB - Key: %s, Value: %v", key, value)

	// Stubbed out "success" - in actual use, you would check for real success/failure from your DB call
	return nil
}

// Global LocalStorage and its mutex
var (
	localStorage      = make(map[string]interface{})
	localStorageMutex sync.RWMutex
)

// SetData sets a key-value pair in the local storage
func SetDataInServerDB(key string, value interface{}) {
	localStorageMutex.Lock()
	defer localStorageMutex.Unlock()

	localStorage[key] = value
}

// GetData retrieves the value for a given key from the local storage
func GetDataFromServerDB(key string) (interface{}, bool) {
	localStorageMutex.RLock()
	defer localStorageMutex.RUnlock()

	value, exists := localStorage[key]
	return value, exists
}

// DeleteData removes a key-value pair from the local storage
func DeleteDataFromServerDB(key string) {
	localStorageMutex.Lock()
	defer localStorageMutex.Unlock()

	delete(localStorage, key)
}
