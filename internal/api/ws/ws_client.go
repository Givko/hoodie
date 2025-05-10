package ws

import (
	"sync"

	"github.com/givko/hoodie/internal/api/ws/connection"
	"github.com/givko/hoodie/internal/api/ws/proto"
	"github.com/go-logr/logr"
)

const (
	MaxConnections = 3
)

type Client struct {
	username      string
	handlers      [MaxConnections]WsHandlerInterface
	handlersMutex sync.RWMutex
	hub           WsHubInterface
	logger        logr.Logger
}

var _ WsClientInterface = (*Client)(nil)

func NewClient(username string, hub WsHubInterface, logger logr.Logger) WsClientInterface {
	return &Client{
		username: username,
		handlers: [MaxConnections]WsHandlerInterface{},
		hub:      hub,
		logger:   logger.WithName("ws_client").WithValues("username", username),
	}
}

// addNewConnection adds a new connection to the client
// It starts the writer and reader goroutines
func (c *Client) AddNewConnection(conn connection.WsConnectionInterface) {
	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()
	for index, handler := range c.handlers {
		if handler != nil {
			continue
		}

		handler := NewWsHandler(conn, c.username, c, c.logger, uint32(index))
		c.handlers[index] = handler
		go handler.Run()
		return
	}

	//Change signature to return error
}

// writeMessage writes a message to all connections of the client
func (c *Client) WriteMessage(message *proto.Message) {
	c.handlersMutex.RLock()
	defer c.handlersMutex.RUnlock()

	for _, handler := range c.handlers {
		if handler == nil {
			continue
		}

		err := handler.WriteMessage(message)
		if err != nil {
			id, _ := handler.GetId()
			c.logger.Error(err, "error writing message", "message", message, "handler", id)
		}
	}
}

// Close closes the provided connection
func (c *Client) Close(conn WsHandlerInterface) error {
	id, _ := conn.GetId()
	err := conn.Close()
	if err != nil {
		c.logger.Error(err, "error closing connection", "username", c.username)
		return err
	}

	c.handlersMutex.Lock()
	defer c.handlersMutex.Unlock()
	c.handlers[id] = nil
	allNil := true
	for _, handler := range c.handlers {
		if handler != nil {
			allNil = false
			break
		}
	}
	if allNil {
		c.hub.Unregister(c)
	}

	return nil
}

// Broadcast sends a message to the central hub
func (c *Client) Broadcast(message *proto.Message) {
	c.hub.Broadcast(message)
}

// GetUsername returns the username of the client
func (c *Client) GetUsername() (string, error) {
	return c.username, nil
}
