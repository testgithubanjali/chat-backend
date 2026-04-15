package websockets

type Hub struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan OutboundMessage
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan OutboundMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.Clients[client.ID] = client

		case client := <-h.Unregister:
			delete(h.Clients, client.ID)
			close(client.Send)

		case msg := <-h.Broadcast:
			for _, client := range h.Clients {

				// ✅ SEND ONLY TO sender + receiver
				if client.ID == msg.ReceiverID || client.ID == msg.SenderID {
					select {
					case client.Send <- msg:
					default:
						close(client.Send)
						delete(h.Clients, client.ID)
					}
				}
			}
		}
	}
}
