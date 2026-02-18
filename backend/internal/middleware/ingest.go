package middleware

import (
	"net/http"
)

// WithIngestSecret returns middleware that validates the X-Ingest-Secret header.
// This is used for server-to-server Garmin ingestion routes (Python sync script).
func WithIngestSecret(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if secret == "" {
				// No secret configured — reject all requests to protect the route.
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			if r.Header.Get("X-Ingest-Secret") != secret {
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
