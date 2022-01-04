package twilio

import (
	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/channels"
	"github.com/caos/zitadel/internal/notification/messages"
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

func (t *Twilio) HandleMessage(message channels.Message) error {
	twilioMsg, ok := message.(*messages.SMS)
	if !ok {
		return caos_errs.ThrowInternal(nil, "TWILI-s0pLc", "message is not SMS")
	}
	m, err := t.client.Messages.SendMessage(twilioMsg.SenderPhoneNumber, twilioMsg.RecipientPhoneNumber, twilioMsg.GetContent(), nil)
	if err != nil {
		return caos_errs.ThrowInternal(err, "TWILI-osk3S", "could not send message")
	}
	logging.LogWithFields("SMS_-f335c523", "message_sid", m.Sid, "status", m.Status).Debug("sms sent")
	return nil
}
