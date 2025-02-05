package ws

import "github.com/givko/hoodie/internal/domain"

type Client struct {
	username string
	handlers []WsHandlerInterface
}

var _ WsClientInterface = (*Client)(nil)

func NewClient(username string) WsClientInterface {
	return &Client{
		username: username,
		handlers: make([]WsHandlerInterface, 0),
	}
}

// addNewConnection adds a new connection to the client
// It starts the writer and reader goroutines
func (c *Client) AddNewConnection(handler WsHandlerInterface) {
	c.handlers = append(c.handlers, handler)

	go handler.Run()
}

// writeMessage writes a message to all connections of the client
func (c *Client) WriteMessage(message *domain.ChatMessage) {
	for _, handler := range c.handlers {
		handler.WriteMessage(message)
	}
}
