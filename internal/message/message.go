package message

type Message struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}

func NewMessage(sender, recipient, content string) *Message {
	return &Message{
		Sender:    sender,
		Recipient: recipient,
		Content:   content,
	}
}
