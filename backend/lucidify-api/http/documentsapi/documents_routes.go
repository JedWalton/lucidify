package documentsapi

import (
	"lucidify-api/server/config"
	"lucidify-api/server/middleware"
	"lucidify-api/service/clerkservice"
	"lucidify-api/service/documentservice"
	"net/http"
)

func SetupRoutes(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) *http.ServeMux {
	mux = SetupDocumentsUploadHandler(config, mux, documentService, clerkService)
	mux = SetupDocumentsGetDocumentHandler(config, mux, documentService, clerkService)
	mux = SetupDocumentsGetAllDocumentHandler(config, mux, documentService, clerkService)
	mux = SetupDocumentsDeleteDocumentHandler(config, mux, documentService, clerkService)
	mux = SetupDocumentsUpdateDocumentNameHandler(config, mux, documentService, clerkService)
	mux = SetupDocumentsUpdateDocumentContentHandler(config, mux, documentService, clerkService)

	return mux
}

func setupHandlerWithMiddleware(config *config.ServerConfig, handler http.Handler, clerkService clerkservice.ClerkClient) http.Handler {
	handlerWithSession := clerkService.WithActiveSession(handler)

	handlerWithSessionFunc := handlerWithSession.(http.HandlerFunc)
	handlerWithSessionFunc = middleware.CORSMiddleware(config.AllowedOrigins)(handlerWithSessionFunc)
	handlerWithSessionFunc = middleware.Logging(handlerWithSessionFunc)

	return handlerWithSessionFunc
}

func SetupDocumentsUploadHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) *http.ServeMux {
	handler := DocumentsUploadHandler(documentService, clerkService)
	mux.Handle("/documents/upload", setupHandlerWithMiddleware(config, handler, clerkService))

	return mux
}

func SetupDocumentsGetDocumentHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) *http.ServeMux {
	handler := DocumentsGetDocumentHandler(documentService, clerkService)
	mux.Handle("/documents/getdocument", setupHandlerWithMiddleware(config, handler, clerkService))

	return mux
}

func SetupDocumentsGetAllDocumentHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) *http.ServeMux {
	handler := DocumentsGetAllDocumentsHandler(documentService, clerkService)
	mux.Handle("/documents/get_all_documents", setupHandlerWithMiddleware(config, handler, clerkService))

	return mux
}

func SetupDocumentsDeleteDocumentHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) *http.ServeMux {
	handler := DocumentsDeleteDocumentHandler(documentService, clerkService)
	mux.Handle("/documents/deletedocument", setupHandlerWithMiddleware(config, handler, clerkService))

	return mux
}

func SetupDocumentsUpdateDocumentNameHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) *http.ServeMux {
	handler := DocumentsUpdateDocumentNameHandler(documentService, clerkService)
	mux.Handle("/documents/update_document_name", setupHandlerWithMiddleware(config, handler, clerkService))

	return mux
}

func SetupDocumentsUpdateDocumentContentHandler(config *config.ServerConfig, mux *http.ServeMux, documentService documentservice.DocumentService, clerkService clerkservice.ClerkClient) *http.ServeMux {
	handler := DocumentsUpdateDocumentContentHandler(documentService, clerkService)
	mux.Handle("/documents/update_document_content", setupHandlerWithMiddleware(config, handler, clerkService))

	return mux
}
