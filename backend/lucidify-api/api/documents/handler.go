package documents

import (
	"encoding/json"
	"log"
	"lucidify-api/modules/store"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func DocumentsUploadHandler(db *store.Store, clerkInstance clerk.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		user, err := clerkInstance.Users().Read(sessClaims.Claims.Subject)
		if err != nil {
			panic(err)
		}

		w.Write([]byte(*&user.ID))

		var reqBody map[string]string
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		document_name := reqBody["document_name"]
		content := reqBody["content"]

		log.Printf("Title: %s\n", document_name)
		log.Printf("Content: %s\n", content)

		// placeholderUserID := "PLACEHOLDER USER ID"
		db.UploadDocument(user.ID, document_name, content)

		// responseBody := map[string]string{
		// 	"user_id": user.ID,
		// }

		w.Header().Set("Content-Type", "application/json")
		// json.NewEncoder(w).Encode(responseBody)
	}
}
