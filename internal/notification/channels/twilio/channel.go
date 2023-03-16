package twilio

import (
	"context"
	"net/url"

	"github.com/kevinburke/twilio-go"
	"github.com/zitadel/logging"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

func InitTwilioChannel(config TwilioConfig) channels.NotificationChannel {
	client := twilio.NewClient(config.SID, config.Token, nil)

	logging.Debug("successfully initialized twilio sms channel")

	return channels.HandleMessageFunc(func(ctx context.Context, message channels.Message) error {
		twilioMsg, ok := message.(*messages.SMS)
		if !ok {
			return caos_errs.ThrowInternal(nil, "TWILI-s0pLc", "message is not SMS")
		}
		v := url.Values{}
		v.Set("Body", twilioMsg.GetContent())
		v.Set("From", twilioMsg.SenderPhoneNumber)
		v.Set("To", twilioMsg.RecipientPhoneNumber)
		m, err := client.Messages.Create(ctx, v)
		if err != nil {
			return caos_errs.ThrowInternal(err, "TWILI-osk3S", "could not send message")
		}
		logging.WithFields("message_sid", m.Sid, "status", m.Status).Debug("sms sent")
		return nil
	})
}
