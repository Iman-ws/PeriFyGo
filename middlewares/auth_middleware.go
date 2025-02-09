package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"PeriFyGo/utils"
)

// JSON-ответ
func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// AuthMiddleware – проверяет JWT токен + роль пользователя
func AuthMiddleware(next http.Handler, requiredRole string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Убираем "Bearer " из заголовка
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			jsonResponse(w, http.StatusUnauthorized, "Unauthorized: Invalid token")
			return
		}

		// Если требуется роль, проверяем её
		if requiredRole != "" && claims.Role != requiredRole {
			jsonResponse(w, http.StatusForbidden, "Forbidden: Access denied")
			return
		}

		// Всё ок – передаем управление дальше
		next.ServeHTTP(w, r)
	})
}
