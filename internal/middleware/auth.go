package middleware

import (
	"context"
	"net/http"
	"strings"

	"prompt-management/internal/config"
	"prompt-management/pkg/auth"
	"prompt-management/pkg/response"
)

type contextKey string

const UserIDContextKey contextKey = "user_id"

// Authenticate is a middleware that validates the JWT in the Authorization header.
func Authenticate(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		token := parts[1]
		claims, err := auth.ValidateToken(token, cfg.JWTSecret)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		// Inject user ID into context
		ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
