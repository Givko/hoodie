package ws

import "github.com/givko/hoodie/internal/domain"

type WsClientInterface interface {
	WriteMessage(message *domain.ChatMessage)
	AddNewConnection(conn WsHandlerInterface)
}
