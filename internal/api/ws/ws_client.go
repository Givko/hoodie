package ws

import "github.com/givko/hoodie/internal/domain"

type Client struct {
	username    string
	connections []*WsConnection
}

func NewClient(username string) *Client {
	return &Client{
		username:    username,
		connections: make([]*WsConnection, 0),
	}
}

// addNewConnection adds a new connection to the client
// It starts the writer and reader goroutines
func (c *Client) addNewConnection(conn *WsConnection) {
	c.connections = append(c.connections, conn)

	go conn.runWriter()
	go conn.runReader()
}

// writeMessage writes a message to all connections of the client
func (c *Client) writeMessage(message domain.ChatMessage) {
	for _, conn := range c.connections {
		conn.writer <- message
	}
}
