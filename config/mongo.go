package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB is a global variable that holds our MongoDB client.
// We will initialize it in InitMongoDB() and use it elsewhere.
var DB *mongo.Client

// InitMongoDB connects to MongoDB using the MONGO_URI environment variable.
func InitMongoDB() *mongo.Client {
	mongoURI := os.Getenv("MONGO_URI") // e.g. "mongodb://localhost:27017"
	if mongoURI == "" {
		log.Fatal("MONGO_URI is not set in environment")
	}

	// Create a new client with the given URI
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Error creating MongoDB client: %v", err)
	}

	// We set a timeout to avoid hanging connections
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to connect to MongoDB
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Cannot connect to MongoDB: %v", err)
	}

	// Just to be sure, let's do a quick Ping
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Cannot ping MongoDB: %v", err)
	}

	log.Println("Successfully connected to MongoDB!")
	DB = client
	return client
}
