package ws

type Hub struct {
	clients    map[string]*Client
	broadcast  chan Message
	Register   chan *WsConnection
	unregister chan *WsConnection
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan Message),
		Register:   make(chan *WsConnection),
		unregister: make(chan *WsConnection),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.registerConn(conn)
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

func (h *Hub) registerConn(conn *WsConnection) {
	if client, ok := h.clients[conn.username]; !ok {
		client = NewClient(conn.username, []*WsConnection{conn})
		h.clients[conn.username] = client
	} else {
		client.addNewConnection(conn)
	}
}

func (h *Hub) broadcastMessage(message Message) {

	client, ok := h.clients[message.Username]
	if !ok {
		//TODO: log error
		return
	}

	client.writeJson(message)
}
