package twilio

import (
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/providers"
	twilio "github.com/kevinburke/twilio-go"
)

type Twilio struct {
	client *twilio.Client
}

func InitTwilioProvider(config TwilioConfig) *Twilio {
	return &Twilio{
		client: twilio.NewClient(config.SID, config.Token, nil),
	}
}

func (t *Twilio) CanHandleMessage(message providers.Message) bool {
	twilioMsg, ok := message.(*TwilioMessage)
	if !ok {
		return false
	}
	return twilioMsg.Content != "" && twilioMsg.RecipientPhoneNumber != "" && twilioMsg.SenderPhoneNumber != ""
}

func (t *Twilio) HandleMessage(message providers.Message) error {
	twilioMsg, ok := message.(*TwilioMessage)
	if !ok {
		return caos_errs.ThrowInternal(nil, "TWILI-s0pLc", "message is not TwilioMessage")
	}
	m, err := t.client.Messages.SendMessage(twilioMsg.SenderPhoneNumber, twilioMsg.RecipientPhoneNumber, twilioMsg.GetContent(), nil)
	if err != nil {
		return caos_errs.ThrowInternal(err, "TWILI-osk3S", "could not send message")
	}
	logging.LogWithFields("SMS_-f335c523", "message_sid", m.Sid, "status", m.Status).Debug("sms sent")
	return nil
}
