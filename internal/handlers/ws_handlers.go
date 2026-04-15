
package handlers

import (
	"chat-backend/internal/models"
	"chat-backend/internal/store"
	"chat-backend/internal/websockets"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RegisterRoutes(e *echo.Echo, hub *websockets.Hub) {

	e.GET("/", func(c echo.Context) error {
		return c.File("web/index.html")
	})

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

	// Get chat history
	e.GET("/messages/:user_id", func(c echo.Context) error {
		userID := c.Param("user_id")

		var msgs []models.Message
		store.DB.Where(
			"sender_id = ? OR receiver_id = ?",
			userID, userID,
		).Order("created_at asc").Find(&msgs)

		return c.JSON(http.StatusOK, msgs)
	})
}