package messages

import (
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/notification/channels"
)

var _ channels.Message = (*SMS)(nil)

type SMS struct {
	SenderPhoneNumber    string
	RecipientPhoneNumber string
	Content              string
	TriggeringEvent      eventstore.Event
}

func (msg *SMS) GetContent() (string, error) {
	return msg.Content, nil
}

func (msg *SMS) GetTriggeringEvent() eventstore.Event {
	return msg.TriggeringEvent
}
