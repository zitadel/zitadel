package twilio

import (
	"github.com/kevinburke/twilio-go"
	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

func InitChannel(config Config) channels.NotificationChannel[*messages.SMS] {
	client := twilio.NewClient(config.SID, config.Token, nil)

	logging.Debug("successfully initialized twilio sms channel")

	return channels.HandleMessageFunc[*messages.SMS](func(message *messages.SMS) error {
		content, err := message.GetContent()
		if err != nil {
			return err
		}
		m, err := client.Messages.SendMessage(message.SenderPhoneNumber, message.RecipientPhoneNumber, content, nil)
		if err != nil {
			return caos_errs.ThrowInternal(err, "TWILI-osk3S", "could not send message")
		}
		logging.WithFields("message_sid", m.Sid, "status", m.Status).Debug("sms sent")
		return nil
	})
}
