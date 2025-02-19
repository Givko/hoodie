package ws

import (
	"fmt"
	"sync"

	"github.com/givko/hoodie/internal/domain"
	"github.com/givko/hoodie/internal/infrastructure/utils"
	"github.com/go-logr/logr"
)

type Client struct {
	username string
	handlers sync.Map
	hub      WsHubInterface
	logger   logr.Logger
}

var _ WsClientInterface = (*Client)(nil)

func NewClient(username string, hub WsHubInterface, logger logr.Logger) WsClientInterface {
	return &Client{
		username: username,
		handlers: sync.Map{},
		hub:      hub,
		logger:   logger,
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
			c.logger.Error(fmt.Errorf("error casting to WsHandlerInterface"), "error casting to WsHandlerInterface", "username", c.username)
			return true
		}

		err := handler.WriteMessage(message)
		if err != nil {
			c.logger.Error(err, "error writing message", "username", c.username, "message", message, "handler", handler)
			return true
		}

		return true
	})
}

// Close closes the provided connection
func (c *Client) Close(conn WsHandlerInterface) error {
	id, _ := conn.GetId()
	err := conn.Close()
	if err != nil {
		c.logger.Error(err, "error closing connection", "username", c.username)
		return err
	}

	c.handlers.Delete(id)
	isEmpty := utils.IsEmpty(&c.handlers)
	if isEmpty {
		fmt.Println("The sync.Map is empty")
	} else {
		fmt.Println("The sync.Map is not empty")
	}

	if isEmpty {
		c.hub.Unregister(c)
	}

	return nil
}

// Broadcast sends a message to the central hub
func (c *Client) Broadcast(message domain.ChatMessage) {
	c.hub.Broadcast(message)
}

// GetUsername returns the username of the client
func (c *Client) GetUsername() (string, error) {
	return c.username, nil
}
