package middleware

import (
	"bytes"
	"io"
	"lucidify-api/modules/config"
	"net/http"

	svix "github.com/svix/svix-webhooks/go"
)

func ClerkWebhooksAuthenticationMiddleware(config *config.ServerConfig) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Read the request body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to process request", http.StatusInternalServerError)
				return
			}

			// Reset the request body so it can be read again by the next handler
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			// Create a new Svix webhook instance using the ClerkSigningSecret from the configuration
			wh, err := svix.NewWebhook(config.ClerkSigningSecret)
			if err != nil {
				http.Error(w, "Failed to process request", http.StatusInternalServerError)
				return
			}

			// Verify the payload using Svix
			err = wh.Verify(body, r.Header)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// If the signature is valid, call the next handler
			next(w, r)
		}
	}
}
