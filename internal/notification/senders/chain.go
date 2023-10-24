package senders

import "github.com/zitadel/zitadel/internal/notification/channels"

var _ channels.NotificationChannel[channels.Message] = (*Chain[channels.Message])(nil)

type Chain[T channels.Message] struct {
	channels []channels.NotificationChannel[T]
}

func ChainChannels[T channels.Message](channel ...channels.NotificationChannel[T]) *Chain[T] {
	return &Chain[T]{channels: channel}
}

// HandleMessage returns a non nil error from a provider immediately if any occurs
// messages are sent to channels in the same order they were provided to ChainChannels()
func (c *Chain[T]) HandleMessage(message T) error {
	for i := range c.channels {
		if err := c.channels[i].HandleMessage(message); err != nil {
			return err
		}
	}
	return nil
}

func (c *Chain[T]) Len() int {
	return len(c.channels)
}
