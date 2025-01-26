package ws

type Client struct {
	username    string
	connections []*WsConnection
}

func NewClient(username string) *Client {
	return &Client{
		username:    username,
		connections: make([]*WsConnection, 0),
	}
}

func (c *Client) addNewConnection(conn *WsConnection) {
	c.connections = append(c.connections, conn)

	go conn.runWriter()
	go conn.runReader()
}

func (c *Client) writeJson(message Message) {
	for _, conn := range c.connections {
		conn.writer <- message
	}
}
