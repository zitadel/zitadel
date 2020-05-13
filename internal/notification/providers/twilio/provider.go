package twilio

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/notification/providers"
	twilio "github.com/kevinburke/twilio-go"
)

type Twilio struct {
	client *twilio.Client
}

func InitTwilioProvider(config *TwilioConfig) (*Twilio, error) {
	return &Twilio{
		client: twilio.NewClient(config.SID, config.Token, nil),
	}, nil
}

func (t *Twilio) CanHandleMessage(message providers.Message) bool {
	twilioMsg := message.(TwilioMessage)
	return twilioMsg.Content != "" && twilioMsg.RecipientPhoneNumber != "" && twilioMsg.SenderPhoneNumber != ""
}

func (t *Twilio) HandleMessage(message providers.Message) error {
	twilioMsg := message.(TwilioMessage)
	m, err := t.client.Messages.SendMessage(twilioMsg.SenderPhoneNumber, twilioMsg.RecipientPhoneNumber, twilioMsg.GetContent(), nil)
	if err != nil {
		return err
	}
	logging.LogWithFields("SMS_-f335c523", "message_sid", m.Sid, "status", m.Status).Debug("sms sent")
	return nil
}
