package ws

import (
	"sync"
	"time"

	"github.com/givko/hoodie/internal/api/ws/connection"
	"github.com/givko/hoodie/internal/domain"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
)

// TODO: make this configurable
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer in bytes
	maxMessageSize = 10 * 1024
)

type WebSocketHandler struct {
	id         string
	username   string
	connection connection.WsConnectionInterface
	writer     chan domain.ChatMessage
	client     WsClientInterface
	once       sync.Once
	logger     logr.Logger
}

var _ WsHandlerInterface = (*WebSocketHandler)(nil)

func NewWsHandler(conn connection.WsConnectionInterface, username string, client WsClientInterface, log logr.Logger) *WebSocketHandler {
	id := uuid.NewString()
	return &WebSocketHandler{
		connection: conn,
		id:         id,
		writer:     make(chan domain.ChatMessage),
		client:     client,
		username:   username,
		logger:     log.WithName("ws_handler").WithValues("id", id, "username", username),
	}
}

func (w *WebSocketHandler) SetClient(client WsClientInterface) {
	w.client = client
}

func (w *WebSocketHandler) GetUsername() (string, error) {
	return w.username, nil
}

func (w *WebSocketHandler) Run() {
	w.logger.Info("Starting websocket handler")
	go w.runReader()
	go w.runWriter()

	w.logger.Info("Websocket handler started")
}

func (w *WebSocketHandler) WriteMessage(message domain.ChatMessage) error {
	w.writer <- message
	return nil
}

// Run starts the connection
// It starts listening for messages from the websocket connection
// and broadcasts them to the hub
func (w *WebSocketHandler) runReader() {
	defer w.client.Close(w)

	w.connection.SetReadLimit(maxMessageSize)
	w.connection.SetReadDeadline(time.Now().Add(pongWait))
	w.connection.SetPongHandler(func(string) error {
		w.connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		message, err := w.connection.ReadMessage()
		if err != nil {
			w.logger.Error(err, "Received message", "chat_message", message)
			break
		}

		w.client.Broadcast(message)
	}
}

// Run starts the writer
// It listens for messages from the writer channel and writes them to the websocket connection
func (w *WebSocketHandler) runWriter() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		w.client.Close(w)
	}()

	for {
		select {
		case message, ok := <-w.writer:
			{
				w.connection.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					w.logger.Info("Connection closed", "connection", w.id, "username", w.username)
					w.connection.WriteCloseMessage()
					return
				} else {
					err := w.connection.WriteMessage(message)
					if err != nil {
						w.logger.Error(err, "Write message unsuccessful", "chat_message", message, "connection", w.id, "username", w.username)
						break
					}
				}
			}
		case <-ticker.C:
			{
				w.connection.SetWriteDeadline(time.Now().Add(writeWait))
				if err := w.connection.WritePingMessage(); err != nil {
					w.logger.Error(err, "Write ping message unsuccessful", "connection", w.id, "username", w.username)
					return
				}
			}
		}
	}
}

func (w *WebSocketHandler) GetId() (string, error) {
	return w.id, nil
}

func (w *WebSocketHandler) Close() error {
	w.once.Do(func() {
		if err := w.connection.Close(); err != nil {
			w.logger.Error(err, "Error closing connection", "connection", w.id, "username", w.username)
		}
	})
	return nil
}
