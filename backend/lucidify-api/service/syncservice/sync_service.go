package syncservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// data := map[string]string{
	// 	"exampleKey": "exampleValue",
	// }
	//
	// response := ServerResponse{
	// 	Success: true,
	// 	Data:    data,
	// 	Message: "Data fetched successfully",
	// }
	fetchedData, err := fetchDataFromDB(key)
	if err != nil {
		response := ServerResponse{
			Success: false,
			Message: err.Error(),
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

	// For demonstration, we'll just print the key and value.
	log.Printf("Syncing data to DB - Key: %s, Value: %v", key, value)

	// Stubbed out "success" - in actual use, you would check for real success/failure from your DB call
	return nil
}

// fetchDataFromDB is a stub function to simulate database fetching.
func fetchDataFromDB(key string) (interface{}, error) {
	// Instead of fetching data from a database, we return a hardcoded value.
	// You should replace this with actual database interaction logic.
	data := map[string]string{
		"apiKey": "exampleValue",
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
