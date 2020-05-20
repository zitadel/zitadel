package chat

type ChatMessage struct {
	Text string `json:"text"`
}

func (msg *ChatMessage) GetContent() string {
	return msg.Text
}
