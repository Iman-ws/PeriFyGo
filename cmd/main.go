package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"PeriFyGo/config"
	"PeriFyGo/routes"

	"github.com/gorilla/handlers"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Init MongoDB
	mongoClient := config.InitMongoDB()

	// Register routes
	r := routes.RegisterRoutes()

	// ✅ Добавляем CORS Middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),                                       // Разрешить все источники (лучше заменить на ngrok)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Разрешить методы
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),           // Разрешить заголовки
	)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler(r), // Используем middleware CORS
	}

	go func() {
		log.Println("🚀 Server running on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Disconnect from Mongo
	if err := mongoClient.Disconnect(ctx); err != nil {
		fmt.Println("Error disconnecting from Mongo:", err)
	}

	log.Println("Server exited properly.")
}
