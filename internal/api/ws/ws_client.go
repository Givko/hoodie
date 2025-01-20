package ws

type Client struct {
	username    string
	connections []*WsConnection
}

func NewClient(username string, conn []*WsConnection) *Client {
	return &Client{
		username:    username,
		connections: conn,
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
