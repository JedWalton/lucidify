package documentsapi

import (
	"encoding/json"
	"fmt"
	"lucidify-api/service/documentservice"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/google/uuid"
)

func DocumentsUploadHandler(documentService documentservice.DocumentService, clerkInstance clerk.Client) http.HandlerFunc {
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

func DocumentsGetDocumentHandler(documentService documentservice.DocumentService, clerkInstance clerk.Client) http.HandlerFunc {
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

func DocumentsGetAllDocumentsHandler(documentService documentservice.DocumentService, clerkInstance clerk.Client) http.HandlerFunc {
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

		document, err := documentService.GetAllDocuments(user.ID)
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

func DocumentsDeleteDocumentHandler(documentService documentservice.DocumentService, clerkInstance clerk.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
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
		documentID := reqBody["documentID"]

		err = documentService.DeleteDocument(user.ID, uuid.MustParse(documentID))
		if err != nil {
			http.Error(w, "Internal server error. Unable to delete document", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func DocumentsUpdateDocumentNameHandler(documentService documentservice.DocumentService, clerkInstance clerk.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
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

		documentID := reqBody["documentID"]
		newDocumentName := reqBody["new_document_name"]

		fmt.Println(documentID)

		err = documentService.UpdateDocumentName(user.ID, uuid.MustParse(documentID), newDocumentName)
		if err != nil {
			http.Error(w, "Internal server error. Unable to update document", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func DocumentsUpdateDocumentContentHandler(documentService documentservice.DocumentService, clerkInstance clerk.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
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

		documentID := reqBody["documentID"]
		newDocumentContent := reqBody["new_document_content"]

		fmt.Println(documentID)

		err = documentService.UpdateDocumentContent(user.ID, uuid.MustParse(documentID), newDocumentContent)
		if err != nil {
			http.Error(w, "Internal server error. Unable to update document", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}
