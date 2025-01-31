package ws

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
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

type WsConnection struct {
	id         string
	username   string
	connection *websocket.Conn
	writer     chan ChatMessage
	hub        *Hub
}

func NewWsConnection(conn *websocket.Conn, hub *Hub, username string) *WsConnection {
	return &WsConnection{
		connection: conn,
		id:         uuid.NewString(),
		writer:     make(chan ChatMessage),
		hub:        hub,
		username:   username,
	}
}

func (w *WsConnection) runReader() {
	defer w.connection.Close()

	w.connection.SetReadLimit(maxMessageSize)
	w.connection.SetReadDeadline(time.Now().Add(pongWait))
	w.connection.SetPongHandler(func(string) error {
		w.connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		message, err := w.readProtobufMessage()
		if err != nil {
			break
		}
		w.hub.broadcast <- message
	}
}

func (w *WsConnection) runWriter() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		w.connection.Close()
	}()

	for {
		select {
		case message, ok := <-w.writer:
			{
				w.connection.SetWriteDeadline(time.Now().Add(writeWait))
				if !ok {

					// The hub closed the channel.
					w.connection.WriteMessage(websocket.CloseMessage, []byte{})
					return
				} else {
					err := w.writeProtobufMessage(message)
					if err != nil {
						fmt.Println("Error: ", err.Error())
						break
					}
				}
			}
		case <-ticker.C:
			{
				w.connection.SetWriteDeadline(time.Now().Add(writeWait))
				if err := w.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}

func (w *WsConnection) readProtobufMessage() (ChatMessage, error) {
	typeId, message, err := w.connection.ReadMessage()
	if err != nil {
		return ChatMessage{}, err
	}

	if typeId != websocket.BinaryMessage {
		return ChatMessage{}, fmt.Errorf("unexpected message type: %d", typeId)
	}

	unmarshaledMessage := Message{}
	err = proto.Unmarshal(message, &unmarshaledMessage)
	if err != nil {
		return ChatMessage{}, err
	}

	return ChatMessage{
		Sender:    unmarshaledMessage.Sender,
		Content:   unmarshaledMessage.Content,
		Recipient: unmarshaledMessage.Recipient,
	}, nil
}

func (w *WsConnection) writeProtobufMessage(message ChatMessage) error {
	marshaledMessage, err := proto.Marshal(&Message{
		Sender:    message.Sender,
		Recipient: message.Recipient,
		Content:   message.Content,
	})
	if err != nil {
		return err
	}

	return w.connection.WriteMessage(websocket.BinaryMessage, marshaledMessage)
}

type ChatMessage struct {
	Sender    string
	Recipient string
	Content   string
}
