package ws

import (
	"sync"

	"github.com/givko/hoodie/internal/domain"
	"github.com/go-logr/logr"
)

type Hub struct {
	clients    sync.Map
	broadcast  chan domain.ChatMessage
	register   chan WsHandlerInterface
	unregister chan WsClientInterface
	logger     logr.Logger
}

var _ WsHubInterface = (*Hub)(nil)

func NewHub(logger logr.Logger) *Hub {
	return &Hub{
		clients:    sync.Map{},
		broadcast:  make(chan domain.ChatMessage, 256),
		register:   make(chan WsHandlerInterface, 256),
		unregister: make(chan WsClientInterface, 256),
		logger:     logger,
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
func (h *Hub) Unregister(conn WsClientInterface) {
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
	newClient := NewClient(username, h, h.logger)

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

func (h *Hub) unregisterConn(conn WsClientInterface) {
	username, err := conn.GetUsername()
	if err != nil {
		h.logger.Error(err, "error getting username from client")
		return
	}

	h.clients.Delete(username)
}

func (h *Hub) getClient(username string) (*Client, bool) {
	value, ok := h.clients.Load(username)
	if !ok {
		return nil, false
	}
	client, ok := value.(*Client)
	return client, ok
}
