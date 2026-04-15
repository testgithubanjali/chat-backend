package main

import (
	"log"
	"os"

	"chat-backend/internal/handlers"
	"chat-backend/internal/models"
	"chat-backend/internal/store"
	"chat-backend/internal/websockets"

	"github.com/labstack/echo/v4"
)

func main() {

	// ENV
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "1234")
	os.Setenv("DB_NAME", "testdb")

	// DB
	store.ConnectDatabase()
	store.DB.AutoMigrate(&models.Message{}) // ✅ FIX

	// Hub
	hub := websockets.NewHub()
	go hub.Run()

	// Server
	e := echo.New()

	handlers.RegisterRoutes(e, hub)

	log.Println("Server running on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}