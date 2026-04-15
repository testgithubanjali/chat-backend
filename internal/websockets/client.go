package websockets

import (
	"context"
	"log"
	"time"

	"chat-backend/internal/models"
	"chat-backend/internal/store"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan OutboundMessage
	Hub  *Hub
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg InboundMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		// ✅ Save message to MongoDB
		collection := store.DB.Collection("messages")

		_, err = collection.InsertOne(context.TODO(), models.Message{
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
			CreatedAt:  time.Now(),
		})

		if err != nil {
			log.Println("Error inserting message:", err)
		}

		// ✅ Broadcast message
		c.Hub.Broadcast <- OutboundMessage{
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
