package channels

import "context"

type Message interface {
	GetContent() string
}

type NotificationChannel interface {
	HandleMessage(context.Context, Message) error
}

var _ NotificationChannel = (HandleMessageFunc)(nil)

type HandleMessageFunc func(context.Context, Message) error

func (h HandleMessageFunc) HandleMessage(ctx context.Context, message Message) error {
	return h(ctx, message)
}
