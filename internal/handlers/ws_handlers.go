package handlers

import (
	"net/http"

	"chat-backend/internal/websockets"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func RegisterRoutes(e *echo.Echo, hub *websockets.Hub) {

	e.GET("/ws", func(c echo.Context) error {

		userID := c.QueryParam("user_id")

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		client := &websockets.Client{
			ID:   userID,
			Conn: conn,
			Send: make(chan websockets.OutboundMessage),
			Hub:  hub,
		}

		hub.Register <- client

		go client.WritePump()
		go client.ReadPump()

		return nil
	})
}
