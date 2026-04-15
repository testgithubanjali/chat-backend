package websockets
import (
	"chat-backend/internal/models"
	"chat-backend/internal/store"
)
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg InboundMessage  // ✅ FIX name
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}

		// Save to DB
		store.DB.Create(&models.Message{
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
		})

		// Broadcast
		c.Hub.Broadcast <- OutboundMessage{
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content:    msg.Content,
		}
	}
}