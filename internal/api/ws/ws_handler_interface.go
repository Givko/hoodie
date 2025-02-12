package ws

import "github.com/givko/hoodie/internal/domain"

type WsHandlerInterface interface {
	WriteMessage(message domain.ChatMessage) error
	Run()
	GetUsername() (string, error)
	GetId() (string, error)
	Close() error
	SetClient(client WsClientInterface)
}
