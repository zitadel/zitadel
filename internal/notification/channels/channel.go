package channels

import "github.com/zitadel/zitadel/v2/internal/eventstore"

type Message interface {
	GetTriggeringEvent() eventstore.Event
	GetContent() (string, error)
}

type NotificationChannel interface {
	HandleMessage(message Message) error
}

var _ NotificationChannel = (HandleMessageFunc)(nil)

type HandleMessageFunc func(message Message) error

func (h HandleMessageFunc) HandleMessage(message Message) error {
	return h(message)
}
