package ws

import (
	"github.com/givko/hoodie/internal/domain"
)

type Hub struct {
	clients    map[string]WsClientInterface
	broadcast  chan *domain.ChatMessage
	register   chan WsHandlerInterface
	unregister chan WsHandlerInterface
}

var _ WsHubInterface = (*Hub)(nil)

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]WsClientInterface),
		broadcast:  make(chan *domain.ChatMessage),
		register:   make(chan WsHandlerInterface),
		unregister: make(chan WsHandlerInterface),
	}
}

// Broadcast sends a message to the recipient
func (h *Hub) Broadcast(message *domain.ChatMessage) {
	h.broadcast <- message
}

// Register registers a new connection
func (h *Hub) Register(conn WsHandlerInterface) {
	h.register <- conn
}

// Unregister unregisters a connection
func (h *Hub) Unregister(conn WsHandlerInterface) {
	h.unregister <- conn
}

// Run starts the hub
// It listens for new connections and messages to broadcast
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.registerConn(conn)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerConn registers a new connection
// It creates a new client if it does not exist and adds the connection to it
// otherwise it adds the connection to the existing client
func (h *Hub) registerConn(conn WsHandlerInterface) {
	username, err := conn.GetUsername()
	if err != nil {
		// TODO log error
		return
	}
	if client, ok := h.clients[username]; !ok {
		client = NewClient(username)
		client.AddNewConnection(conn)
		h.clients[username] = client
	} else {
		client.AddNewConnection(conn)
	}
}

// broadcastMessage broadcasts a message to the recipient
// It finds the client by the recipient username and sends the message to all connections of the client
func (h *Hub) broadcastMessage(message *domain.ChatMessage) {

	client, ok := h.clients[message.Recipient]
	if !ok {
		//TODO: log error
		return
	}

	client.WriteMessage(message)
}
