package websockets
import (
	"chat-backend/internal/models"
	"chat-backend/internal/store"
)

type Client struct {
	ID   string
	Conn *Connection
	Send chan OutboundMessage
	Hub *Hub

}
func (c *Client) ReadPump() {
	defer func() {	
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		var msg InbpoundMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}
      store.DB.Create(&models.Message{
			SenderID: msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content: msg.Content,
		})
		c.Hub.Broadcast <- OutboundMessage{
			SenderID: msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Content: msg.Content,
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