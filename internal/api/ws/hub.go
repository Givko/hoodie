package ws

import (
	"context"
	"sync"

	"github.com/givko/hoodie/internal/api/ws/connection"
	"github.com/givko/hoodie/internal/api/ws/proto"
	"github.com/go-logr/logr"
	"github.com/redis/go-redis/v9"
)

type Hub struct {
	clients     sync.Map
	broadcast   chan *proto.Message
	redisClient *redis.Client
	register    chan RegisterPair
	unregister  chan WsClientInterface
	logger      logr.Logger
}

type RegisterPair struct {
	Username string
	Conn     connection.WsConnectionInterface
}

var _ WsHubInterface = (*Hub)(nil)

func NewHub(logger logr.Logger, redisClient *redis.Client) *Hub {
	return &Hub{
		clients:     sync.Map{},
		broadcast:   make(chan *proto.Message, 256),
		register:    make(chan RegisterPair, 256),
		unregister:  make(chan WsClientInterface, 256),
		logger:      logger.WithName("ws_hub"),
		redisClient: redisClient,
	}
}

// Broadcast sends a message to the recipient
func (h *Hub) Broadcast(message *proto.Message) {
	h.broadcast <- message
}

// Register registers a new connection
func (h *Hub) Register(conn RegisterPair) {
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
func (h *Hub) registerConn(registerPair RegisterPair) {
	h.logger.Info("registering new connection", "username", registerPair.Username)
	username := registerPair.Username

	// Create a new client candidate.
	newClient := NewClient(username, h, h.logger)

	// Atomically store or retrieve the client.
	actual, _ := h.clients.LoadOrStore(username, newClient)
	client := actual.(*Client)

	client.AddNewConnection(registerPair.Conn)
	h.logger.Info("registered new connection", "username", registerPair.Username)
}

// broadcastMessage broadcasts a message to the recipient
// It finds the client by the recipient username and sends the message to all connections of the client
func (h *Hub) broadcastMessage(message *proto.Message) {
	ctx := context.Background()
	err := h.redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: "chat_messages",
		Values: map[string]interface{}{
			"content":   message.Content,
			"sender":    message.Sender,
			"recipient": message.Recipient,
		}}).Err()

	if err != nil {
		h.logger.Error(err, "error publishing message to redis", "recipient", message.Recipient)
		return
	}

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
