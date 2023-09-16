package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Read the request body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		// Log the request body
		log.Printf("Request Body: %s", bodyBytes)

		// Reset the request body to its original state
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("Received %s request for %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		log.Println("Handled the request.")
	}
}
