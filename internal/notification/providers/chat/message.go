package chat

type ChatMessage struct {
	Content string
}

func (msg *ChatMessage) GetContent() string {
	return msg.Content
}
