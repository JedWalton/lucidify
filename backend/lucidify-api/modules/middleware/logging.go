package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// func Logging(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Read the request body
// 		bodyBytes, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Printf("Error reading body: %v", err)
// 			http.Error(w, "can't read body", http.StatusBadRequest)
// 			return
// 		}
//
// 		// Log the request body
// 		log.Printf("Request Body: %s", bodyBytes)
//
// 		// Log the headers
// 		for name, values := range r.Header {
// 			for _, value := range values {
// 				log.Printf("Header: %s = %s", name, value)
// 			}
// 		}
//
// 		// Reset the request body to its original state
// 		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
//
// 		log.Printf("Received %s request for %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
// 		next(w, r)
// 		log.Println("Handled the request.")
// 	}
// }

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

		// Log the headers in a structured manner
		var headerStrings []string
		for name, values := range r.Header {
			headerStrings = append(headerStrings, fmt.Sprintf("%s: %s", name, strings.Join(values, ", ")))
		}
		log.Printf("Headers: {\n\t%s\n}", strings.Join(headerStrings, "\n\t"))

		// Reset the request body to its original state
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Printf("Received %s request for %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next(w, r)
		log.Println("Handled the request.")
	}
}
