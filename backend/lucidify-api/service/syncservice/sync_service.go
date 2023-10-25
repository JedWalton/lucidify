package syncservice

import "fmt"

// ServerResponse is the structure that defines the standard response from the server.
type ServerResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Holds the actual data, if any
	Message string      `json:"message,omitempty"` // Descriptive message, especially useful in case of errors
}

// This is a utility function to send JSON responses

// func IsValidKey(key string) bool {
// 	switch key {
// 	case "apiKey", "ConversationHistory", "SelectedConversation", "Theme", "Folders",
// 		"Prompts", "ShowChatbar", "ShowPromptbar", "PluginKeys", "Settings":
// 		return true
// 	}
// 	return false
// }

// func HandleFetch(key string) (interface{}, ServerResponse)
// func HandleDelete(key string) ServerResponse
// func HandleSync(key string, value interface{}) ServerResponse

func HandleSet(key string, value interface{}) ServerResponse {
	ok := SetDataInLocalStorage(key, value)
	if !ok {
		return ServerResponse{Success: false, Message: "error setting data"}
	}
	return ServerResponse{Success: true, Message: "Data synced successfully"}
}

func HandleGet(key string) (interface{}, ServerResponse) {
	data, ok := GetDataFromLocalStorage(key)
	if ok && data != "" {
		return data, ServerResponse{Success: true, Message: "Data fetched successfully"}
	}
	return nil, ServerResponse{Success: false, Message: "No data found for key: " + key}
}

func HandleRemove(key string) ServerResponse {
	ok := RemoveDataFromLocalStorage(key)
	if ok {
		return ServerResponse{Success: true, Message: "Data deleted successfully"}
	}
	return ServerResponse{Success: false, Message: "Data not deleted. unsuccessful"}
}

// SyncDataToDb function that accepts a key and value and syncs this data with a database.
// You should replace the contents of this function with actual database interaction logic.
func SyncData(key string, value interface{}) error {
	// This function is currently a stub and does not actually interact with a database.
	// Here, you would write your logic to sync data to your database.
	// This might include SQL statements, or calls to another service, etc.

	// SetDataInServerDB(key, value)
	// For demonstration, we'll just print the key and value.
	// log.Printf("Syncing data to DB - Key: %s, Value: %v", key, value)
	ok := SetDataInLocalStorage(key, value)
	if !ok {
		return fmt.Errorf("invalid key. Not set in localstorage: %s", key)
	}

	// Stubbed out "success" - in actual use, you would check for real success/failure from your DB call
	return nil
}

var Storage LocalStorage // Global variable representing our local storage

// GetDataFromLocalStorage retrieves a value from local storage based on key.
func GetDataFromLocalStorage(key string) (interface{}, bool) {
	switch key {
	case "apiKey":
		return Storage.APIKey, true
	case "conversationHistory":
		return Storage.ConversationHistory, true
	case "selectedConversation":
		return Storage.SelectedConversation, true
	case "theme":
		return Storage.Theme, true
	case "folders":
		return Storage.Folders, true
	case "prompts":
		return Storage.Prompts, true
	case "showChatbar":
		return Storage.ShowChatbar, true
	case "showPromptbar":
		return Storage.ShowPromptbar, true
	case "pluginKeys":
		return Storage.PluginKeys, true
	case "settings":
		return Storage.Settings, true
	default:
		return nil, false
	}
}

// SetDataInLocalStorage sets a value in local storage based on key.
func SetDataInLocalStorage(key string, value interface{}) bool {
	switch key {
	case "apiKey":
		Storage.APIKey = value.(string)
	case "conversationHistory":
		Storage.ConversationHistory = value.([]Conversation)
	case "selectedConversation":
		Storage.SelectedConversation = value.(Conversation)
	case "theme":
		Storage.Theme = value.(string)
	case "folders":
		Storage.Folders = value.([]FolderInterface)
	case "prompts":
		Storage.Prompts = value.([]Prompt)
	case "showChatbar":
		Storage.ShowChatbar = value.(bool)
	case "showPromptbar":
		Storage.ShowPromptbar = value.(bool)
	case "pluginKeys":
		Storage.PluginKeys = value.([]PluginKey)
	case "settings":
		Storage.Settings = value.(Settings)
	default:
		return false
	}
	return true
}

// RemoveDataFromLocalStorage removes a value from local storage based on key.
func RemoveDataFromLocalStorage(key string) bool {
	switch key {
	case "apiKey":
		Storage.APIKey = ""
		return true
	case "conversationHistory":
		Storage.ConversationHistory = nil
		return true
	case "selectedConversation":
		Storage.SelectedConversation = Conversation{}
		return true
	case "theme":
		Storage.Theme = ""
		return true
	case "folders":
		Storage.Folders = nil
		return true
	case "prompts":
		Storage.Prompts = nil
		return true
	case "showChatbar":
		Storage.ShowChatbar = false
		return true
	case "showPromptbar":
		Storage.ShowPromptbar = false
		return true
	case "pluginKeys":
		Storage.PluginKeys = nil
		return true
	case "settings":
		Storage.Settings = Settings{}
		return true
	default:
		return false
	}
}
