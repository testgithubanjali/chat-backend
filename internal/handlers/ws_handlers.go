package handlers

import (
	"context"
	"net/http"

	"chat-backend/internal/models"
	"chat-backend/internal/store"
	"chat-backend/internal/websockets"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RegisterRoutes(e *echo.Echo, hub *websockets.Hub) {

	// ✅ Health check
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is running 🚀")
	})

	// ✅ Serve frontend (optional)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Chat Backend Running")
	})

	// ✅ WebSocket endpoint
	e.GET("/ws", func(c echo.Context) error {

		userID := c.QueryParam("user_id")

		if userID == "" {
			return c.String(http.StatusBadRequest, "user_id required")
		}

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

	// ✅ Get chat history
	e.GET("/messages/:user_id", func(c echo.Context) error {

		userID := c.Param("user_id")

		collection := store.DB.Collection("messages")

		filter := bson.M{
			"$or": []bson.M{
				{"sender_id": userID},
				{"receiver_id": userID},
			},
		}

		cursor, err := collection.Find(context.TODO(), filter)
		if err != nil {
			return err
		}
		defer cursor.Close(context.TODO())

		var msgs []models.Message

		for cursor.Next(context.TODO()) {
			var msg models.Message
			if err := cursor.Decode(&msg); err != nil {
				return err
			}
			msgs = append(msgs, msg)
		}

		return c.JSON(http.StatusOK, msgs)
	})
}
