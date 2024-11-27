package twilio

import (
	"errors"

	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func InitChannel(config Config) channels.NotificationChannel {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{Username: config.SID, Password: config.Token})
	logging.Debug("successfully initialized twilio sms channel")

	return channels.HandleMessageFunc(func(message channels.Message) error {
		twilioMsg, ok := message.(*messages.SMS)
		if !ok {
			return zerrors.ThrowInternal(nil, "TWILI-s0pLc", "message is not SMS")
		}
		if config.VerifyServiceSID != "" {
			params := &verify.CreateVerificationParams{}
			params.SetTo(twilioMsg.RecipientPhoneNumber)
			params.SetChannel("sms")

			resp, err := client.VerifyV2.CreateVerification(config.VerifyServiceSID, params)

			var twilioErr *twilioClient.TwilioRestError
			if errors.As(err, &twilioErr) && twilioErr.Code == 60203 {
				// If there were too many attempts to send a verification code (more than 5 times)
				// without a verification check, even retries with backoff might not solve the problem.
				// Instead, let the user initiate the verification again (e.g. using "resend code")
				// https://www.twilio.com/docs/api/errors/60203
				logging.WithFields("error", twilioErr.Message, "code", twilioErr.Code).Warn("twilio create verification error")
				return channels.NewCancelError(twilioErr)
			}

			if err != nil {
				return zerrors.ThrowInternal(err, "TWILI-0s9f2", "could not send verification")
			}
			logging.WithFields("sid", resp.Sid, "status", resp.Status).Debug("verification sent")

			twilioMsg.VerificationID = resp.Sid
			return nil
		}

		content, err := twilioMsg.GetContent()
		if err != nil {
			return err
		}
		params := &openapi.CreateMessageParams{}
		params.SetTo(twilioMsg.RecipientPhoneNumber)
		params.SetFrom(twilioMsg.SenderPhoneNumber)
		params.SetBody(content)
		m, err := client.Api.CreateMessage(params)
		if err != nil {
			return zerrors.ThrowInternal(err, "TWILI-osk3S", "could not send message")
		}
		logging.WithFields("message_sid", m.Sid, "status", m.Status).Debug("sms sent")
		return nil
	})
}
