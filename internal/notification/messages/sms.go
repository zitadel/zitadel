package messages

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
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

var _ channels.Message = (*TwilioVerify)(nil)

type TwilioVerify struct {
	VerifyServiceSID     string
	RecipientPhoneNumber string
	TriggeringEvent      eventstore.Event
}

func (msg *TwilioVerify) GetContent() (string, error) {
	return "", nil
}
func (msg *TwilioVerify) GetTriggeringEvent() eventstore.Event {
	return msg.TriggeringEvent
}
