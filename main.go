package main

import (
	"log"
	"net/http"

	"chat-backend/internal/handlers"
	"chat-backend/internal/websocket"
)

func main() {
	hub := websocket.NewHub()

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleWebSocket(hub, w, r)
	})

	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}