package syncapi

import (
	"encoding/json"
	"log"
	"lucidify-api/service/syncservice"
	"net/http"
)

// ServerResponse is the structure that defines the standard response from the server.
type ServerResponse struct {
	Success bool        `json:"success"`           // Indicates if the operation was successful
	Data    interface{} `json:"data,omitempty"`    // Holds the actual data, if any
	Message string      `json:"message,omitempty"` // Descriptive message, especially useful in case of errors
}

func SyncHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Set content type for all responses from this handler

		log.Printf("Request method: %s, URL: %s", r.Method, r.URL.String())

		switch r.Method {
		case http.MethodGet, http.MethodDelete, http.MethodPost:
			// For GET, DELETE, and POST, read 'key' from query parameters
			key := r.URL.Query().Get("key")
			if key == "" {
				response := ServerResponse{
					Success: false,
					Message: "Key not provided",
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(response)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			if r.Method == http.MethodGet {
				syncservice.FetchData(w, r, key)
			} else if r.Method == http.MethodDelete {
				syncservice.DeleteData(w, r, key)
			} else {
				// For POST, read 'value' from the request body
				var requestData map[string]interface{}
				err := json.NewDecoder(r.Body).Decode(&requestData)
				if err != nil {
					http.Error(w, "Bad request data", http.StatusBadRequest)
					return
				}

				value, valueExists := requestData["value"]
				if !valueExists {
					http.Error(w, "Value not provided in request body", http.StatusBadRequest)
					return
				}

				// Now you have the 'key' from the URL and 'value' from the request body and can proceed
				// You might want to modify your 'syncDataToDB' function to accept both 'key' and 'value'
				err = syncservice.SyncData(key, value) // Make sure this function accepts both key and value
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// Return a successful response if there were no errors
				response := ServerResponse{
					Success: true,
					Message: "Data synced successfully",
				}
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(response)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
