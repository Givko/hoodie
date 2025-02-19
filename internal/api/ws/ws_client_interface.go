package ws

import (
	"github.com/givko/hoodie/internal/api/ws/connection"
	"github.com/givko/hoodie/internal/domain"
)

type WsClientInterface interface {
	WriteMessage(message domain.ChatMessage)
	AddNewConnection(conn connection.WsConnectionInterface)
	Close(conn WsHandlerInterface) error
	Broadcast(message domain.ChatMessage)
	GetUsername() (string, error)
}
