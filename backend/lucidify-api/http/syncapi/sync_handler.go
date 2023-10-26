package syncapi

import (
	"encoding/json"
	"io"
	"log"
	"lucidify-api/service/syncservice"
	"net/http"
)

// This is a utility function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, statusCode int, response syncservice.ServerResponse) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SyncHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		key := r.URL.Query().Get("key")

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			// Handle the error, maybe return a response indicating the error.
			return
		}
		value := string(bodyBytes)

		var response syncservice.ServerResponse

		switch r.Method {
		case http.MethodGet:
			resp := syncservice.HandleGet(key)
			response = resp
			if response.Success {
				response.Data = resp.Data
			}
		case http.MethodDelete:
			response = syncservice.HandleRemove(key)
		case http.MethodPost:
			response = syncservice.HandleSet(key, value)
		default:
			response = syncservice.ServerResponse{
				Success: false,
				Message: "Method not allowed",
			}
			sendJSONResponse(w, http.StatusMethodNotAllowed, response)
			return
		}
		sendJSONResponse(w, http.StatusOK, response)
	}
}
