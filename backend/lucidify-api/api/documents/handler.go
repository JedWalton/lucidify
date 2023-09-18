package documents

import (
	"encoding/json"
	"log"
	"lucidify-api/modules/store"
	"net/http"
)

func DocumentsUploadHandler(store *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqBody map[string]string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		document_name := reqBody["document_name"]
		content := reqBody["content"]

		log.Printf("Title: %s\n", document_name)
		log.Printf("Content: %s\n", content)

		placeholderUserID := 1
		store.UploadDocument(placeholderUserID, document_name, content)

		responseMessage := "PLACEHOLDER RESPONSE"

		responseBody := map[string]string{
			"response": responseMessage,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseBody)
	}
}
