package ws

import (
	"fmt"
	"sync"
	"time"

	"github.com/givko/hoodie/internal/api/ws/connection"
	"github.com/givko/hoodie/internal/domain"
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
}

var _ WsHandlerInterface = (*WebSocketHandler)(nil)

func NewWsHandler(conn connection.WsConnectionInterface, username string) *WebSocketHandler {
	return &WebSocketHandler{
		connection: conn,
		id:         uuid.NewString(),
		writer:     make(chan domain.ChatMessage),
		client:     nil,
		username:   username,
	}
}

func (w *WebSocketHandler) SetClient(client WsClientInterface) {
	w.client = client
}

func (w *WebSocketHandler) GetUsername() (string, error) {
	return w.username, nil
}

func (w *WebSocketHandler) Run() {
	go w.runReader()
	go w.runWriter()
}

func (w *WebSocketHandler) WriteMessage(message domain.ChatMessage) error {
	w.writer <- message
	return nil
}

// Run starts the connection
// It starts listening for messages from the websocket connection
// and broadcasts them to the hub
func (w *WebSocketHandler) runReader() {
	defer w.Close()

	w.connection.SetReadLimit(maxMessageSize)
	w.connection.SetReadDeadline(time.Now().Add(pongWait))
	w.connection.SetPongHandler(func(string) error {
		w.connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		message, err := w.connection.ReadMessage()
		if err != nil {
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
		w.Close()
	}()

	for {
		select {
		case message, ok := <-w.writer:
			{
				w.connection.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {
					w.connection.WriteCloseMessage()
					w.client.Close(w)
					return
				} else {
					err := w.WriteMessage(message)
					if err != nil {
						fmt.Println("Error: ", err.Error())
						break
					}
				}
			}
		case <-ticker.C:
			{
				w.connection.SetWriteDeadline(time.Now().Add(writeWait))
				if err := w.connection.WritePingMessage(); err != nil {
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
	f := sync.OnceFunc(func() { w.connection.Close() })
	f()
	return nil
}
