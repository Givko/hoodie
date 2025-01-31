package ws

import "github.com/givko/hoodie/internal/domain"

type Hub struct {
	clients    map[string]*Client
	broadcast  chan domain.ChatMessage
	Register   chan *WsConnection
	unregister chan *WsConnection
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan domain.ChatMessage),
		Register:   make(chan *WsConnection),
		unregister: make(chan *WsConnection),
	}
}

// Run starts the hub
// It listens for new connections and messages to broadcast
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.registerConn(conn)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerConn registers a new connection
// It creates a new client if it does not exist and adds the connection to it
// otherwise it adds the connection to the existing client
func (h *Hub) registerConn(conn *WsConnection) {
	if client, ok := h.clients[conn.username]; !ok {
		client = NewClient(conn.username)
		client.addNewConnection(conn)
		h.clients[conn.username] = client
	} else {
		client.addNewConnection(conn)
	}
}

// broadcastMessage broadcasts a message to the recipient
// It finds the client by the recipient username and sends the message to all connections of the client
func (h *Hub) broadcastMessage(message domain.ChatMessage) {

	client, ok := h.clients[message.Recipient]
	if !ok {
		//TODO: log error
		return
	}

	client.writeMessage(message)
}
