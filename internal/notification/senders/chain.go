package senders

import "github.com/zitadel/zitadel/internal/notification/channels"

var _ channels.NotificationChannel = (*Chain)(nil)

type Chain struct {
	channels []channels.NotificationChannel
}

func ChainChannels(channel ...channels.NotificationChannel) *Chain {
	return &Chain{channels: channel}
}

// HandleMessage returns a non nil error from a provider immediately if any occurs
// messages are sent to channels in the same order they were provided to ChainChannels()
func (c *Chain) HandleMessage(message channels.Message) error {
	for i := range c.channels {
		if err := c.channels[i].HandleMessage(message); err != nil {
			return err
		}
	}
	return nil
}

func (c *Chain) Len() int {
	return len(c.channels)
}
