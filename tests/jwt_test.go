package tests

import (
	"os"
	"testing"
	"time"

	"PeriFyGo/utils"
)

func TestJWTGenerationAndValidation(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret_key")

	// Задаем тестовые данные.
	email := "test@example.com"
	role := "user"

	// Генерируем токен.
	token, err := utils.GenerateToken(email, role)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Проверяем токен.
	claims, err := utils.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Проверяем, что возвращенные claims содержат ожидаемые email и role.
	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("Expected role %s, got %s", role, claims.Role)
	}
}
