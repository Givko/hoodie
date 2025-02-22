package connection

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	wsProto "github.com/givko/hoodie/internal/api/ws/proto"
)

type WsConnectionInterface interface {
	ReadMessage() (*wsProto.Message, error)
	WriteMessage(message *wsProto.Message) error
	WriteCloseMessage() error
	WritePingMessage() error
	Close() error
	SetPongHandler(handler func(string) error)
	SetReadDeadline(t time.Time) error
	SetReadLimit(limit int64)
	SetWriteDeadline(t time.Time) error
}

var _ WsConnectionInterface = (*ConnectionWrapper)(nil)

type ConnectionWrapper struct {
	conn *websocket.Conn
}

func NewConnectionWrapper(conn *websocket.Conn) WsConnectionInterface {
	return &ConnectionWrapper{
		conn: conn,
	}
}

func (c *ConnectionWrapper) ReadMessage() (*wsProto.Message, error) {
	return c.readProtobufMessage()
}

func (c *ConnectionWrapper) WriteMessage(message *wsProto.Message) error {
	return c.writeProtobufMessage(message)
}

func (c *ConnectionWrapper) WriteCloseMessage() error {
	return c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *ConnectionWrapper) Close() error {
	return c.conn.Close()
}

func (c *ConnectionWrapper) SetPongHandler(handler func(string) error) {
	c.conn.SetPongHandler(handler)
}

func (c *ConnectionWrapper) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *ConnectionWrapper) SetReadLimit(limit int64) {
	c.conn.SetReadLimit(limit)
}

func (c *ConnectionWrapper) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func (c *ConnectionWrapper) WritePingMessage() error {
	return c.conn.WriteMessage(websocket.PingMessage, []byte{})
}

// readProtobufMessage reads a protobuf message from the websocket connection
// It returns the message and an error if any
// It returns an error if the message type is not binary
// It returns an error if the message cannot be unmarshaled
func (w *ConnectionWrapper) readProtobufMessage() (*wsProto.Message, error) {
	typeId, message, err := w.conn.ReadMessage()
	if err != nil {
		return &wsProto.Message{}, err
	}

	if typeId != websocket.BinaryMessage {
		return &wsProto.Message{}, fmt.Errorf("unexpected message type: %d", typeId)
	}

	unmarshaledMessage := wsProto.Message{}
	err = proto.Unmarshal(message, &unmarshaledMessage)
	if err != nil {
		return &wsProto.Message{}, err
	}

	return &unmarshaledMessage, nil
}

// writeProtobufMessage writes a protobuf message to the websocket connection
// It returns an error if the message cannot be marshaled
// It returns an error if the message cannot be written to the connection
func (w *ConnectionWrapper) writeProtobufMessage(message *wsProto.Message) error {
	marshaledMessage, err := proto.Marshal(&wsProto.Message{
		Sender:    message.Sender,
		Recipient: message.Recipient,
		Content:   message.Content,
	})
	if err != nil {
		return err
	}

	return w.conn.WriteMessage(websocket.BinaryMessage, marshaledMessage)
}
