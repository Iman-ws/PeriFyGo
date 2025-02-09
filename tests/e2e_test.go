package tests

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func TestLoginAndProfileRedirect(t *testing.T) {
	// Настройка возможностей для headless Chrome.
	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		// Args: []string{"--headless", "--no-sandbox", "--disable-dev-shm-usage"},
		Args: []string{"--no-sandbox", "--disable-dev-shm-usage"},
	}
	caps.AddChrome(chromeCaps)

	// Используем новый URL без /wd/hub для ChromeDriver (Selenium 4+).
	wd, err := selenium.NewRemote(caps, "http://localhost:9515")
	if err != nil {
		t.Fatalf("Failed to create Selenium remote session: %v", err)
	}
	defer wd.Quit()

	// Переходим на страницу входа.
	if err := wd.Get("http://localhost:8080/login.html"); err != nil {
		t.Fatalf("Failed to load login page: %v", err)
	}

	// Находим элементы формы входа.
	emailInput, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {
		t.Fatalf("Email input not found: %v", err)
	}
	passwordInput, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {
		t.Fatalf("Password input not found: %v", err)
	}

	// Вводим данные тестового пользователя.
	if err := emailInput.SendKeys("meralieviman6@gmail.com"); err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
	if err := passwordInput.SendKeys("Iman2006"); err != nil {
		t.Fatalf("Failed to send password: %v", err)
	}

	// Нажимаем кнопку входа.
	loginButton, err := wd.FindElement(selenium.ByCSSSelector, "button[type='submit']")
	if err != nil {
		t.Fatalf("Login button not found: %v", err)
	}
	if err := loginButton.Click(); err != nil {
		t.Fatalf("Failed to click login button: %v", err)
	}

	// Ждем, чтобы произошел редирект (например, 2 секунды).
	time.Sleep(2 * time.Second)

	// Проверяем, что текущий URL равен "http://localhost:8080/profile.html".
	currentURL, err := wd.CurrentURL()
	if err != nil {
		t.Fatalf("Failed to get current URL: %v", err)
	}
	expectedURL := "http://localhost:8080/profile.html"
	if currentURL != expectedURL {
		t.Errorf("Expected URL: %s, got: %s", expectedURL, currentURL)
	}
}
