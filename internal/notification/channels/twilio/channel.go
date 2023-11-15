package twilio

import (
	"github.com/kevinburke/twilio-go"
	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/v2/internal/errors"
	"github.com/zitadel/zitadel/v2/internal/notification/channels"
	"github.com/zitadel/zitadel/v2/internal/notification/messages"
)

func InitChannel(config Config) channels.NotificationChannel {
	client := twilio.NewClient(config.SID, config.Token, nil)

	logging.Debug("successfully initialized twilio sms channel")

	return channels.HandleMessageFunc(func(message channels.Message) error {
		twilioMsg, ok := message.(*messages.SMS)
		if !ok {
			return caos_errs.ThrowInternal(nil, "TWILI-s0pLc", "message is not SMS")
		}
		content, err := twilioMsg.GetContent()
		if err != nil {
			return err
		}
		m, err := client.Messages.SendMessage(twilioMsg.SenderPhoneNumber, twilioMsg.RecipientPhoneNumber, content, nil)
		if err != nil {
			return caos_errs.ThrowInternal(err, "TWILI-osk3S", "could not send message")
		}
		logging.WithFields("message_sid", m.Sid, "status", m.Status).Debug("sms sent")
		return nil
	})
}
