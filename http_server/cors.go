package http_server

import (
	"github.com/go-chi/cors"
	"net/http"
)

// Cors - For more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
func Cors(allowedOrigins []string) func(handler http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}
