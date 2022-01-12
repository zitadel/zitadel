package channels

type Message interface {
	GetContent() string
}

type NotificationChannel interface {
	HandleMessage(message Message) error
}

var _ NotificationChannel = (HandleMessageFunc)(nil)

type HandleMessageFunc func(message Message) error

func (h HandleMessageFunc) HandleMessage(message Message) error {
	return h(message)
}
