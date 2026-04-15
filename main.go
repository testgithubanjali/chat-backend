package main

import (
	"log"

	"chat-backend/internal/handlers"
	"chat-backend/internal/store"
	"chat-backend/internal/websockets"

	"github.com/labstack/echo/v4"
)

func main() {

	// Connect MongoDB
	store.ConnectDatabase()

	// Start WebSocket Hub
	hub := websockets.NewHub()
	go hub.Run()

	// Start server
	e := echo.New()

	handlers.RegisterRoutes(e, hub)

	log.Println("Server running on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
