package ws

import "github.com/givko/hoodie/internal/domain"

type WsHubInterface interface {
	Broadcast(message *domain.ChatMessage)
	Register(conn WsHandlerInterface)
	Unregister(conn WsHandlerInterface)
	Run()
}
