package documentsapi

import (
	"encoding/json"
	"lucidify-api/modules/store/store"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func DocumentsUploadHandler(documentService store.DocumentService, clerkInstance clerk.Client) http.HandlerFunc {
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

		_, err = documentService.UploadDocument(user.ID, document_name, content)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func DocumentsGetDocumentHandler(documentService store.DocumentService, clerkInstance clerk.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
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

		var reqBody map[string]string
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		document_name := reqBody["document_name"]

		document, err := documentService.GetDocument(user.ID, document_name)
		if err != nil {
			http.Error(w, "Internal server error. Unable to get document", http.StatusInternalServerError)
			return
		}

		// Set the Content-Type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Encode the document as JSON and write it to the response writer
		encoder := json.NewEncoder(w)
		err = encoder.Encode(document)
		if err != nil {
			http.Error(w, "Internal server error. Unable to encode document as JSON", http.StatusInternalServerError)
			return
		}
	}
}

//
// func DocumentsGetAllDocumentsHandler(db *postgresqlclient.PostgreSQL, clerkInstance clerk.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}
//
// 		ctx := r.Context()
//
// 		sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
// 		if !ok {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte("Unauthorized"))
// 			return
// 		}
//
// 		user, err := clerkInstance.Users().Read(sessClaims.Claims.Subject)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		document, err := db.GetAllDocuments(user.ID)
// 		if err != nil {
// 			http.Error(w, "Internal server error. Unable to get document", http.StatusInternalServerError)
// 			return
// 		}
//
// 		// Set the Content-Type to application/json
// 		w.Header().Set("Content-Type", "application/json")
//
// 		// Encode the document as JSON and write it to the response writer
// 		encoder := json.NewEncoder(w)
// 		err = encoder.Encode(document)
// 		if err != nil {
// 			http.Error(w, "Internal server error. Unable to encode document as JSON", http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }
//
// func DocumentsDeleteDocumentHandler(db *postgresqlclient.PostgreSQL, clerkInstance clerk.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodDelete {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}
//
// 		ctx := r.Context()
//
// 		sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
// 		if !ok {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte("Unauthorized"))
// 			return
// 		}
//
// 		user, err := clerkInstance.Users().Read(sessClaims.Claims.Subject)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		var reqBody map[string]string
// 		decoder := json.NewDecoder(r.Body)
// 		err = decoder.Decode(&reqBody)
// 		if err != nil {
// 			http.Error(w, "Bad request", http.StatusBadRequest)
// 			return
// 		}
// 		document_name := reqBody["document_name"]
//
// 		err = db.DeleteDocument(user.ID, document_name)
// 		if err != nil {
// 			http.Error(w, "Internal server error. Unable to delete document", http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.Header().Set("Content-Type", "application/json")
// 	}
// }
//
// func DocumentsUpdateDocumentHandler(db *postgresqlclient.PostgreSQL, clerkInstance clerk.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodPut {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}
//
// 		ctx := r.Context()
//
// 		sessClaims, ok := ctx.Value(clerk.ActiveSessionClaims).(*clerk.SessionClaims)
// 		if !ok {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte("Unauthorized"))
// 			return
// 		}
//
// 		user, err := clerkInstance.Users().Read(sessClaims.Claims.Subject)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		var reqBody map[string]string
// 		decoder := json.NewDecoder(r.Body)
// 		err = decoder.Decode(&reqBody)
// 		if err != nil {
// 			http.Error(w, "Bad request", http.StatusBadRequest)
// 			return
// 		}
// 		document_name := reqBody["document_name"]
// 		content := reqBody["content"]
//
// 		err = db.UpdateDocument(user.ID, document_name, content)
// 		if err != nil {
// 			http.Error(w, "Internal server error. Unable to update document", http.StatusInternalServerError)
// 			return
// 		}
//
// 		w.Header().Set("Content-Type", "application/json")
// 	}
// }
