package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
	// MongoDB URI (change if needed)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	// Check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB not reachable:", err)
	}

	log.Println("✅ Connected to MongoDB")

	// Select DB
	DB = client.Database("chat_app")
}