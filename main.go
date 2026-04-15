package main

import (
	"log"

	"chat-backend/internal/handlers"
	"chat-backend/internal/store"
	"chat-backend/internal/websockets"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	store.ConnectDatabase()

	hub := websockets.NewHub()
	go hub.Run()

	e := echo.New()

	// ✅ CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3001",
			"http://localhost",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
	}))

	handlers.RegisterRoutes(e, hub)

	log.Println("Server running on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
