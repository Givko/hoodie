package ws

import (
	"sync"

	"github.com/givko/hoodie/internal/domain"
)

type Hub struct {
	clients    sync.Map
	broadcast  chan domain.ChatMessage
	register   chan WsHandlerInterface
	unregister chan WsHandlerInterface
}

var _ WsHubInterface = (*Hub)(nil)

func NewHub() *Hub {
	return &Hub{
		clients:    sync.Map{},
		broadcast:  make(chan domain.ChatMessage),
		register:   make(chan WsHandlerInterface),
		unregister: make(chan WsHandlerInterface),
	}
}

// Broadcast sends a message to the recipient
func (h *Hub) Broadcast(message domain.ChatMessage) {
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
		case conn := <-h.unregister:
			h.unregisterConn(conn)
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

	// Create a new client candidate.
	newClient := NewClient(username, h)

	// Atomically store or retrieve the client.
	actual, _ := h.clients.LoadOrStore(username, newClient)
	client := actual.(*Client)

	client.AddNewConnection(conn)
}

// broadcastMessage broadcasts a message to the recipient
// It finds the client by the recipient username and sends the message to all connections of the client
func (h *Hub) broadcastMessage(message domain.ChatMessage) {
	client, ok := h.getClient(message.Recipient)
	if !ok {
		return
	}

	client.WriteMessage(message)
}

func (h *Hub) unregisterConn(conn WsHandlerInterface) {
	username, err := conn.GetUsername()
	if err != nil {
		// TODO log error
		return
	}

	if client, ok := h.getClient(username); ok {
		client.Close(conn)
	}
}

func (h *Hub) getClient(username string) (*Client, bool) {
	value, ok := h.clients.Load(username)
	if !ok {
		return nil, false
	}
	client := value.(*Client)
	return client, true
}
