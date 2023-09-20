package middleware

import (
	"bytes"
	"context"
	"io"
	"lucidify-api/modules/config"
	"net/http"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	svix "github.com/svix/svix-webhooks/go"
)

func ClerkAuthenticationMiddleware(config *config.ServerConfig) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Wrap the next handler with Clerk's middleware
			protectedHandler := clerk.WithSessionV2(config.ClerkClient)(next)
			protectedHandler.ServeHTTP(w, r)

			// Retrieve the authenticated session's claims
			session, ok := clerk.SessionFromContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Access the "Subject" field for user ID from the jwt.Claims
			userID := session.Claims.Subject
			if userID == "" {
				http.Error(w, "User ID not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

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
