package middlewares

import (
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authentication
		next.ServeHTTP(w, r)
		// Logging
	})	
}