package ws

type Hub struct {
	clients    map[string]*Client
	broadcast  chan ChatMessage
	Register   chan *WsConnection
	unregister chan *WsConnection
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		broadcast:  make(chan ChatMessage),
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
		client = NewClient(conn.username)
		client.addNewConnection(conn)
		h.clients[conn.username] = client
	} else {
		client.addNewConnection(conn)
	}
}

func (h *Hub) broadcastMessage(message ChatMessage) {

	client, ok := h.clients[message.Recipient]
	if !ok {
		//TODO: log error
		return
	}

	client.writeJson(message)
}
