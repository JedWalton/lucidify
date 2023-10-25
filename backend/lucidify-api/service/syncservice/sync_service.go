package syncservice

import (
	"encoding/json"
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

	// fetchedData, exists := GetDataFromServerDB(key)
	// if !exists {
	// 	response := ServerResponse{
	// 		Success: false,
	// 		Message: "No data found for key: " + key,
	// 	}
	// 	sendJSONResponse(w, http.StatusInternalServerError, response)
	// 	return
	// }
	serverResponse := ServerResponse{
		Success: true,
		Data:    "hello",
		Message: "Data fetched successfully",
	}

	sendJSONResponse(w, http.StatusOK, serverResponse)
}

func DeleteData(w http.ResponseWriter, r *http.Request, key string) {
	response := ServerResponse{
		Success: true,
		Data:    "deleted placeholder",
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

	// SetDataInServerDB(key, value)
	// For demonstration, we'll just print the key and value.
	log.Printf("Syncing data to DB - Key: %s, Value: %v", key, value)

	// Stubbed out "success" - in actual use, you would check for real success/failure from your DB call
	return nil
}

var Storage LocalStorage // Global variable representing our local storage

// GetDataFromLocalStorage retrieves a value from local storage based on key.
func GetDataFromLocalStorage(key string) (interface{}, bool) {
	switch key {
	case "APIKey":
		return Storage.APIKey, true
	case "ConversationHistory":
		return Storage.ConversationHistory, true
	case "SelectedConversation":
		return Storage.SelectedConversation, true
	case "Theme":
		return Storage.Theme, true
	case "Folders":
		return Storage.Folders, true
	case "Prompts":
		return Storage.Prompts, true
	case "ShowChatbar":
		return Storage.ShowChatbar, true
	case "ShowPromptbar":
		return Storage.ShowPromptbar, true
	case "PluginKeys":
		return Storage.PluginKeys, true
	case "Settings":
		return Storage.Settings, true
	case "CHANGELOG":
		return Storage.CHANGELOG, true
	default:
		return nil, false
	}
}

// SetDataInLocalStorage sets a value in local storage based on key.
func SetDataInLocalStorage(key string, value interface{}) bool {
	switch key {
	case "APIKey":
		Storage.APIKey = value.(string)
	case "ConversationHistory":
		Storage.ConversationHistory = value.([]Conversation)
	case "SelectedConversation":
		Storage.SelectedConversation = value.(Conversation)
	case "Theme":
		Storage.Theme = value.(string)
	case "Folders":
		Storage.Folders = value.([]FolderInterface)
	case "Prompts":
		Storage.Prompts = value.([]Prompt)
	case "ShowChatbar":
		Storage.ShowChatbar = value.(bool)
	case "ShowPromptbar":
		Storage.ShowPromptbar = value.(bool)
	case "PluginKeys":
		Storage.PluginKeys = value.([]PluginKey)
	case "Settings":
		Storage.Settings = value.(Settings)
	case "CHANGELOG":
		Storage.CHANGELOG = value.(*[]ChangeLog)
	default:
		return false
	}
	return true
}

// RemoveDataFromLocalStorage removes a value from local storage based on key.
func RemoveDataFromLocalStorage(key string) bool {
	switch key {
	case "APIKey":
		Storage.APIKey = ""
		return true
	case "ConversationHistory":
		Storage.ConversationHistory = nil
		return true
	case "SelectedConversation":
		Storage.SelectedConversation = Conversation{}
		return true
	case "Theme":
		Storage.Theme = ""
		return true
	case "Folders":
		Storage.Folders = nil
		return true
	case "Prompts":
		Storage.Prompts = nil
		return true
	case "ShowChatbar":
		Storage.ShowChatbar = false
		return true
	case "ShowPromptbar":
		Storage.ShowPromptbar = false
		return true
	case "PluginKeys":
		Storage.PluginKeys = nil
		return true
	case "Settings":
		Storage.Settings = Settings{}
		return true
	case "CHANGELOG":
		Storage.CHANGELOG = nil
		return true
	default:
		return false
	}
}
