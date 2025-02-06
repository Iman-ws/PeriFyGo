package middlewares

import (
	"net/http"
	"strings"

	"PeriFyGo/utils"
)

// AuthMiddleware validates JWT token and checks if the user has the required role.
// If requiredRole is an empty string, only token validity is checked.
func AuthMiddleware(next http.Handler, requiredRole string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header. Expected format: "Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token by removing "Bearer " prefix.
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		if tokenString == "" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		// Validate token using our utility function.
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// If a specific role is required, check it.
		if requiredRole != "" && claims.Role != requiredRole {
			http.Error(w, "Access denied: insufficient permissions", http.StatusForbidden)
			return
		}

		// If everything is ok, proceed to the next handler.
		next.ServeHTTP(w, r)
	})
}
