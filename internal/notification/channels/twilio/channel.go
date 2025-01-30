package twilio

import (
	"errors"

	"github.com/twilio/twilio-go"
	twilioClient "github.com/twilio/twilio-go/client"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	aggregateTypeNotification = "notification"
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

			// In case of any client error (4xx), we should not retry sending the verification code
			// as it would be a waste of resources and could potentially result in a rate limit.
			var twilioErr *twilioClient.TwilioRestError
			if errors.As(err, &twilioErr) && twilioErr.Status >= 400 && twilioErr.Status < 500 {
				userID, notificationID := userAndNotificationIDsFromEvent(twilioMsg.TriggeringEvent)
				logging.WithFields(
					"error", twilioErr.Message,
					"status", twilioErr.Status,
					"code", twilioErr.Code,
					"instanceID", twilioMsg.TriggeringEvent.Aggregate().InstanceID,
					"userID", userID,
					"notificationID", notificationID).
					Warn("twilio create verification error")
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

func userAndNotificationIDsFromEvent(event eventstore.Event) (userID, notificationID string) {
	aggID := event.Aggregate().ID

	// we cannot cast to the actual event type because of circular dependencies
	// so we just check the type...
	if event.Aggregate().Type != aggregateTypeNotification {
		// in case it's not a notification event, we can directly return the aggregate ID (as it's a user event)
		return aggID, ""
	}
	// ...and unmarshal the event data from the notification event into a struct that contains the fields we need
	var data struct {
		Request struct {
			UserID string `json:"userID"`
		} `json:"request"`
	}
	if err := event.Unmarshal(&data); err != nil {
		return "", aggID
	}
	return data.Request.UserID, aggID
}
