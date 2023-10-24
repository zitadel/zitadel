package channels

import "github.com/zitadel/zitadel/internal/eventstore"

type Message interface {
	GetTriggeringEvent() eventstore.Event
	GetContent() (string, error)
}

type NotificationChannel[T Message] interface {
	HandleMessage(message T) error
}

var _ NotificationChannel[Message] = (HandleMessageFunc[Message])(nil)

type HandleMessageFunc[T Message] func(message T) error

func (h HandleMessageFunc[T]) HandleMessage(message T) error {
	return h(message)
}
