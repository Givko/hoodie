package ws

import (
	"sync"
	"time"

	"github.com/givko/hoodie/internal/api/ws/connection"
	"github.com/givko/hoodie/internal/api/ws/proto"
	"github.com/go-logr/logr"
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
	indexId    uint32
	connection connection.WsConnectionInterface
	writer     chan *proto.Message
	client     WsClientInterface
	once       sync.Once
	logger     logr.Logger
}

var _ WsHandlerInterface = (*WebSocketHandler)(nil)

func NewWsHandler(
	conn connection.WsConnectionInterface,
	client WsClientInterface,
	log logr.Logger,
	indexId uint32) *WebSocketHandler {
	return &WebSocketHandler{
		connection: conn,
		indexId:    indexId,
		writer:     make(chan *proto.Message, 256),
		client:     client,
		logger:     log.WithName("ws_handler"),
	}
}

func (w *WebSocketHandler) SetClient(client WsClientInterface) {
	w.client = client
}

func (w *WebSocketHandler) Run() {
	w.logger.Info("Starting websocket handler")
	go w.runReader()
	go w.runWriter()

	w.logger.Info("Websocket handler started")
}

func (w *WebSocketHandler) WriteMessage(message *proto.Message) error {
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
					w.logger.Info("Connection closed")
					w.connection.WriteCloseMessage()
					return
				} else {

					err := w.connection.WriteMessage(message)
					if err != nil {
						w.logger.Error(err, "Write message unsuccessful", "chat_message", message)
					}
				}
			}
		case <-ticker.C:
			{
				w.connection.SetWriteDeadline(time.Now().Add(writeWait))
				if err := w.connection.WritePingMessage(); err != nil {
					w.logger.Error(err, "Write ping message unsuccessful")
					return
				}
			}
		}
	}
}

func (w *WebSocketHandler) GetId() (uint32, error) {
	return w.indexId, nil
}

func (w *WebSocketHandler) Close() error {
	w.once.Do(func() {
		if err := w.connection.Close(); err != nil {
			w.logger.Error(err, "Error closing connection")
		}
	})
	return nil
}
