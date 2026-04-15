package websockets
import (
	"chat-backend/internal/models"
	"chat-backend/internal/store"
)
type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan OutboundMessage
	Register   chan *Client//chantype
	Unregister chan *Client
}	
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan OutboundMessage),
		Register:   make(chan *Client),//chan definition
		Unregister: make(chan *Client),
	}	
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
		case message := <-h.Broadcast:
			for client := range h.Clients {
				if client.ID == message.ReceiverID || client.ID == message.SenderID {
					select {
					case client.Send <- message:		
					default:
						close(client.Send)
						delete(h.Clients, client)
					}		
				}
			}
		}
	}
}

