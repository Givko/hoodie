package ws

import "github.com/givko/hoodie/internal/api/ws/proto"

type WsHubInterface interface {
	Broadcast(message *proto.Message)
	Register(conn RegisterPair)
	Unregister(conn WsClientInterface)
	Run()
}
