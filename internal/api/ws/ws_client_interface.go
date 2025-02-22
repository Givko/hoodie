package ws

import (
	"github.com/givko/hoodie/internal/api/ws/connection"
	"github.com/givko/hoodie/internal/api/ws/proto"
)

type WsClientInterface interface {
	WriteMessage(message *proto.Message)
	AddNewConnection(conn connection.WsConnectionInterface)
	Close(conn WsHandlerInterface) error
	Broadcast(message *proto.Message)
	GetUsername() (string, error)
}
