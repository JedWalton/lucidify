package middleware

import "net/http"

func Chain(f http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
