package documentsapi

import (
	"encoding/json"
	"fmt"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
	"net/http"

	"github.com/google/uuid"
)

func DocumentsUploadHandler(documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := clerkService.GetUserIDFromSession(r.Context())
		if err != nil {
			panic(err)
		}

		w.Write([]byte(userID))

		var reqBody map[string]string
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		document_name := reqBody["document_name"]
		content := reqBody["content"]

		_, err = documentService.UploadDocument(userID, document_name, content)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func DocumentsGetDocumentHandler(documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := clerkService.GetUserIDFromSession(r.Context())
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

		document, err := documentService.GetDocument(userID, document_name)
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

func DocumentsGetAllDocumentsHandler(documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := clerkService.GetUserIDFromSession(r.Context())
		if err != nil {
			panic(err)
		}

		document, err := documentService.GetAllDocuments(userID)
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

func DocumentsDeleteDocumentHandler(documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := clerkService.GetUserIDFromSession(r.Context())
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

		err = documentService.DeleteDocument(userID, uuid.MustParse(documentID))
		if err != nil {
			http.Error(w, "Internal server error. Unable to delete document", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func DocumentsUpdateDocumentNameHandler(documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := clerkService.GetUserIDFromSession(r.Context())
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

		err = documentService.UpdateDocumentName(userID, uuid.MustParse(documentID), newDocumentName)
		if err != nil {
			http.Error(w, "Internal server error. Unable to update document", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}

func DocumentsUpdateDocumentContentHandler(documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userID, err := clerkService.GetUserIDFromSession(r.Context())
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

		err = documentService.UpdateDocumentContent(userID, uuid.MustParse(documentID), newDocumentContent)
		if err != nil {
			http.Error(w, "Internal server error. Unable to update document", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}
