package messages

import (
	"fmt"

	twilioclient "github.com/twilio/twilio-go/client"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
)

var _ channels.Message = (*SMS)(nil)

type SMS struct {
	SenderPhoneNumber    string
	RecipientPhoneNumber string
	Content              string
	TriggeringEvent      eventstore.Event

	// VerificationID is set by the sender
	VerificationID *string
}

func (msg *SMS) GetContent() (string, error) {
	return msg.Content, nil
}

func (msg *SMS) GetTriggeringEvent() eventstore.Event {
	return msg.TriggeringEvent
}

type ErrTwilioSendNotification struct {
	Msg         string
	TwilioError twilioclient.TwilioRestError
	Fatal       bool
}

func (e ErrTwilioSendNotification) Error() string {
	return e.Msg
}

func NewTwilioError(twilioError twilioclient.TwilioRestError, fatal bool) error {
	return &ErrTwilioSendNotification{
		Msg:         fmt.Sprintf("%s. Code: %v Status: %v", twilioError.Message, twilioError.Code, twilioError.Status),
		TwilioError: twilioError,
		Fatal:       fatal,
	}
}
