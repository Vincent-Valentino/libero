package middleware

import (
	"context"
	"net/http"
	"strings"
	"libero-backend/internal/service" // Import service package
)

// userContextKey is the key used to store user information in the request context
type userContextKey string

// UserKey is the context key for user data
const UserKey userContextKey = "user"

// AuthMiddleware creates a middleware for JWT authentication
// It now depends on AuthService to validate tokens
func AuthMiddleware(authService service.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate token using AuthService
			claims, err := authService.ValidateJWTToken(tokenString)
			if err != nil {
				// Check for specific error types if needed (e.g., service.ErrTokenInvalid)
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Store validated claims (*service.JWTClaims) in request context
			ctx := context.WithValue(r.Context(), UserKey, claims)

			// Call the next handler with the updated context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware creates a middleware to check user roles
func RoleMiddleware(allowedRoles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user claims (*service.JWTClaims) from context
			claims, ok := r.Context().Value(UserKey).(*service.JWTClaims) // Use service.JWTClaims
			if !ok {
				// This should ideally not happen if AuthMiddleware ran successfully
				http.Error(w, "Unauthorized: User claims not found in context", http.StatusUnauthorized)
				return
			}

			// Check if user has one of the allowed roles
			hasRole := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext extracts user claims (*service.JWTClaims) from the request context
func GetUserFromContext(ctx context.Context) (*service.JWTClaims, bool) { // Return service.JWTClaims
	claims, ok := ctx.Value(UserKey).(*service.JWTClaims) // Use service.JWTClaims
	return claims, ok
}