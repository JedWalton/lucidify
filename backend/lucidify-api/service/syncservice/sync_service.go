package syncservice

// ServerResponse is the structure that defines the standard response from the server.
type ServerResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Holds the actual data, if any
	Message string      `json:"message,omitempty"` // Descriptive message, especially useful in case of errors
}

var store = make(map[string]string)

func HandleSet(key string, value string) ServerResponse {
	store[key] = value
	return ServerResponse{Success: true, Message: "Data set successfully for key: " + key}
}

func HandleGet(key string) ServerResponse {
	if data, ok := store[key]; ok {
		return ServerResponse{Success: true, Data: data, Message: "Data fetched successfully"}
	}
	return ServerResponse{Success: false, Message: "Data not found for key: " + key}
}

// func HandleGet(key string) ServerResponse {
// 	// data, ok := GetDataFromLocalStorage(key)
// 	// if ok && data != "" {
// 	// 	return data, ServerResponse{Success: true, Message: "Data fetched successfully"}
// 	// }
// 	return ServerResponse{Success: true, Message: "Successful Get placeholder for key: " + key}
// }
//
// func HandleSet(key string, value string) ServerResponse {
// 	// ok := SetDataInLocalStorage(key, value)
// 	// if !ok {
// 	// 	return ServerResponse{Success: false, Message: "error setting data"}
// 	// }
// 	return ServerResponse{Success: true, Message: "Successful Set placeholder for key: " + key + " with value: " + value}
// }
//
// func HandleRemove(key string) ServerResponse {
// 	// ok := RemoveDataFromLocalStorage(key)
// 	// if ok {
// 	// 	return ServerResponse{Success: true, Message: "Data deleted successfully"}
// 	// }
// 	return ServerResponse{Success: true, Message: "Successful Delete placeholder for key: " + key}
// }

func HandleRemove(key string) ServerResponse {
	if _, ok := store[key]; ok {
		delete(store, key)
		return ServerResponse{Success: true, Message: "Data deleted successfully for key: " + key}
	}
	return ServerResponse{Success: false, Message: "No data to delete for key: " + key}
}
