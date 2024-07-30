package twilio

import (
	"github.com/kevinburke/twilio-go"
	newTwilio "github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func InitChannel(config Config) channels.NotificationChannel {
	client := newTwilio.NewRestClientWithParams(newTwilio.ClientParams{Username: config.SID, Password: config.Token})
	logging.Debug("successfully initialized twilio sms channel")

	return channels.HandleMessageFunc(func(message channels.Message) error {
		if twilioMsg, ok := message.(*messages.TwilioVerify); ok {

			// SEE: https://www.twilio.com/docs/verify/api/verification
			params := &verify.CreateVerificationParams{}
			params.SetTo(twilioMsg.RecipientPhoneNumber)
			params.SetChannel("sms")
		
			resp, err := client.VerifyV2.CreateVerification(config.VerifyServiceSID, params)
			if err != nil {
				return zerrors.ThrowInternal(err, "TWILI-0s9f2", "could not send verification")
			}


			// How user code verification should happen (Response status will be "approved" if correct code is provided)
			// SEE: https://www.twilio.com/docs/verify/api/verification-check
			// checkParams := &verify.CreateVerificationCheckParams{}
			// checkParams.SetTo(twilioMsg.RecipientPhoneNumber)
			// checkParams.SetCode(userCode) // This will be the code the user entered
			// client.VerifyV2.CreateVerificationCheck(config.VerifyServiceSID, checkParams)
			logging.WithFields("sid", resp.Sid, "status", resp.Status).Debug("verification sent")
			return nil
		}

		client := twilio.NewClient(config.SID, config.Token, nil)
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
