package syncapi

import (
	"encoding/json"
	"lucidify-api/service/syncservice"
	"net/http"
)

func SyncHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json") // Set content type for all responses from this handler

		switch r.Method {
		case http.MethodGet:
			// Logic for fetching data by key, replace with actual logic
			data, err := syncservice.FetchDataFromDB("exampleKey") // this function should be implemented
			if err != nil {
				http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
				return
			}

			response := map[string]interface{}{
				"status": "success",
				"data":   data,
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)

		case http.MethodPost:
			// Logic for saving/updating data
			syncservice.SyncData(w, r) // this function should be implemented and return an error if it fails
			// err := SyncData(w, r) // this function should be implemented and return an error if it fails
			// if err != nil {
			// 	http.Error(w, "Failed to sync data", http.StatusInternalServerError)
			// 	return
			// }

			response := map[string]string{
				"status":  "success",
				"message": "Data synced successfully",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)

		case http.MethodDelete:
			// Logic for deleting data by key
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(w, "Key not provided", http.StatusBadRequest)
				return
			}

			syncservice.DeleteData(w, r, key) // this function should be implemented and return an error if it fails
			// err := DeleteData(w, r, key) // this function should be implemented and return an error if it fails
			// if err != nil {
			// 	http.Error(w, "Failed to delete data", http.StatusInternalServerError)
			// 	return
			// }

			response := map[string]string{
				"status":  "success",
				"message": "Data deleted successfully",
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
