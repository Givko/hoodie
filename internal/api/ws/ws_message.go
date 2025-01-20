package ws

//Message struct to hold message data
type Message struct {
	Username  string `json:"username"`
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	ID        string `json:"id"`
}
