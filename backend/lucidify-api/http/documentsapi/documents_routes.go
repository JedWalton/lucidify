package documentsapi

import (
	"lucidify-api/server/config"
	"lucidify-api/server/middleware"
	"lucidify-api/service/documentservice"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, client clerk.Client) *http.ServeMux {
	mux = SetupDocumentsUploadHandler(config, mux, documentService, client)
	mux = SetupDocumentsGetDocumentHandler(config, mux, documentService, client)
	mux = SetupDocumentsGetAllDocumentHandler(config, mux, documentService, client)
	mux = SetupDocumentsDeleteDocumentHandler(config, mux, documentService, client)
	mux = SetupDocumentsUpdateDocumentNameHandler(config, mux, documentService, client)
	mux = SetupDocumentsUpdateDocumentContentHandler(config, mux, documentService, client)

	return mux
}

func SetupDocumentsUploadHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, client clerk.Client) *http.ServeMux {

	handler := DocumentsUploadHandler(documentService, client)

	injectActiveSession := clerk.WithSession(client)

	handler = middleware.Logging(handler)

	mux.Handle("/documents/upload", injectActiveSession(handler))

	return mux
}

func SetupDocumentsGetDocumentHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, client clerk.Client) *http.ServeMux {

	handler := DocumentsGetDocumentHandler(documentService, client)

	injectActiveSession := clerk.WithSession(client)

	handler = middleware.Logging(handler)

	mux.Handle("/documents/getdocument", injectActiveSession(handler))

	return mux
}

func SetupDocumentsGetAllDocumentHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, client clerk.Client) *http.ServeMux {

	handler := DocumentsGetAllDocumentsHandler(documentService, client)

	injectActiveSession := clerk.WithSession(client)

	handler = middleware.Logging(handler)

	mux.Handle("/documents/get_all_documents", injectActiveSession(handler))

	return mux
}

func SetupDocumentsDeleteDocumentHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, client clerk.Client) *http.ServeMux {

	handler := DocumentsDeleteDocumentHandler(documentService, client)

	injectActiveSession := clerk.WithSession(client)

	handler = middleware.Logging(handler)

	mux.Handle("/documents/deletedocument", injectActiveSession(handler))

	return mux
}

func SetupDocumentsUpdateDocumentNameHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, client clerk.Client) *http.ServeMux {

	handler := DocumentsUpdateDocumentNameHandler(documentService, client)

	injectActiveSession := clerk.WithSession(client)

	handler = middleware.Logging(handler)

	mux.Handle("/documents/update_document_name", injectActiveSession(handler))

	return mux
}

func SetupDocumentsUpdateDocumentContentHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, client clerk.Client) *http.ServeMux {

	handler := DocumentsUpdateDocumentContentHandler(documentService, client)

	injectActiveSession := clerk.WithSession(client)

	handler = middleware.Logging(handler)

	mux.Handle("/documents/update_document_content", injectActiveSession(handler))

	return mux
}
