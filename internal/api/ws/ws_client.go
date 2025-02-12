package ws

import (
	"fmt"
	"sync"

	"github.com/givko/hoodie/internal/domain"
)

type Client struct {
	username string
	handlers sync.Map
	hub      WsHubInterface
}

var _ WsClientInterface = (*Client)(nil)

func NewClient(username string, hub WsHubInterface) WsClientInterface {
	return &Client{
		username: username,
		handlers: sync.Map{},
		hub:      hub,
	}
}

// addNewConnection adds a new connection to the client
// It starts the writer and reader goroutines
func (c *Client) AddNewConnection(handler WsHandlerInterface) {
	id, _ := handler.GetId()
	handler.SetClient(c)
	c.handlers.Store(id, handler)

	go handler.Run()
}

// writeMessage writes a message to all connections of the client
func (c *Client) WriteMessage(message domain.ChatMessage) {

	c.handlers.Range(func(key, value interface{}) bool {
		handler, ok := value.(WsHandlerInterface)
		if !ok {
			return false
		}

		handler.WriteMessage(message)
		return true
	})
}

// Close closes the provided connection
func (c *Client) Close(conn WsHandlerInterface) error {
	id, _ := conn.GetId()
	c.handlers.Delete(id)
	err := conn.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}

	return nil
}

// Broadcast sends a message to the central hub
func (c *Client) Broadcast(message domain.ChatMessage) {
	c.hub.Broadcast(message)
}
