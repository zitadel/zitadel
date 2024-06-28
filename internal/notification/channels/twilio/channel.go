package twilio

import (
	"context"
	"net/url"

	"github.com/kevinburke/twilio-go"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func InitChannel(config Config) channels.NotificationChannel {
	client := twilio.NewClient(config.SID, config.Token, nil)
	logging.Debug("successfully initialized twilio sms channel")

	return channels.HandleMessageFunc(func(message channels.Message) error {
		if twilioMsg, ok := message.(*messages.TwilioVerify); ok {

			// SEE: https://www.twilio.com/docs/verify/api/verification
			body := url.Values{}
			body.Add("To", twilioMsg.RecipientPhoneNumber)
			body.Add("Channel", "sms")
			// NB: by passing the code used by zitadel, we override the default code generation
			// in twilio and therefore can (possibly) skip the seoncdary API call to verify the code
			body.Add("CustomCode", twilioMsg.Code)
			// TODO: pass through custom message (but not code)
			// body.Add("CustomMessage", twilioMsg.Message)

			// TODO: need to pass in parent context
			resp, err := client.Verifications.Create(
				context.Background(),
				twilioMsg.VerifyServiceSID,
				body,
			)
			if err != nil {
				return zerrors.ThrowInternal(err, "TWILI-0s9f2", "could not send verification")
			}
			logging.WithFields("sid", resp.Sid, "status", resp.Status).Debug("verification sent")
			return nil
		}

		twilioMsg, ok := message.(*messages.SMS)
		if !ok {
			return zerrors.ThrowInternal(nil, "TWILI-s0pLc", "message is not SMS")
		}
		content, err := twilioMsg.GetContent()
		if err != nil {
			return err
		}
		m, err := client.Messages.SendMessage(twilioMsg.SenderPhoneNumber, twilioMsg.RecipientPhoneNumber, content, nil)
		if err != nil {
			return zerrors.ThrowInternal(err, "TWILI-osk3S", "could not send message")
		}
		logging.WithFields("message_sid", m.Sid, "status", m.Status).Debug("sms sent")
		return nil
	})
}
