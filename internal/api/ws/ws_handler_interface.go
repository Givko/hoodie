package ws

import "github.com/givko/hoodie/internal/api/ws/proto"

type WsHandlerInterface interface {
	WriteMessage(message *proto.Message) error
	Run()
	GetUsername() (string, error)
	GetId() (string, error)
	Close() error
	SetClient(client WsClientInterface)
}
