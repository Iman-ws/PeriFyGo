package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"PeriFyGo/config"
	"PeriFyGo/controllers"
	"github.com/joho/godotenv"
)

func TestRegisterIntegration(t *testing.T) {
	// Вывод текущей рабочей директории.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Working Directory:", wd)
	t.Log("Absolute path of .env:", filepath.Join(wd, ".env"))

	// Загружаем .env из родительской директории
	err = godotenv.Load("../.env")
	if err != nil {
		t.Log("No .env file found in parent directory, using system environment variables")
	} else {
		t.Log(".env loaded successfully from parent directory")
	}

	// Проверяем, что MONGO_URI установлен.
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		t.Fatal("MONGO_URI is not set in environment")
	} else {
		t.Log("MONGO_URI from test:", mongoURI)
	}

	// Инициализируем подключение к MongoDB.
	config.DB = config.InitMongoDB()
	if config.DB == nil {
		t.Fatal("Failed to initialize MongoDB connection")
	}

	// Подготовка запроса.
	jsonStr := []byte(`{"name":"IntegrationUser","email":"integration_unique@test.com","password":"Test1234"}`)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	authCtrl := &controllers.AuthController{}
	authCtrl.Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d. Response body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}
}
