package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/satishthakur/health-assistant/backend/internal/auth"
)

// contextKey is a private type to avoid context key collisions.
type contextKey string

const userIDKey contextKey = "userID"

// WithAuth returns middleware that validates Bearer JWT tokens.
// On success the userID is injected into the request context.
// On failure a 401 response is returned immediately.
func WithAuth(tokenService *auth.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			userID, err := tokenService.ValidateToken(parts[1])
			if err != nil {
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext retrieves the authenticated user's ID from the context.
// Returns an empty string if not present.
func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey).(string)
	return v
}
