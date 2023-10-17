package types

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/query"
)

func generateSms(
	ctx context.Context,
	channels ChannelChains,
	user *query.NotifyUser,
	content string,
	lastPhone bool,
	triggeringEvent eventstore.Event,
) error {
	number := ""
	smsChannels, twilioConfig, err := channels.SMS(ctx)
	logging.OnError(err).Error("could not create sms channel")
	if smsChannels == nil || smsChannels.Len() == 0 {
		return errors.ThrowPreconditionFailed(nil, "PHONE-w8nfow", "Errors.Notification.Channels.NotPresent")
	}
	if err == nil {
		number = twilioConfig.SenderNumber
	}
	message := &messages.SMS{
		SenderPhoneNumber:    number,
		RecipientPhoneNumber: user.VerifiedPhone,
		Content:              content,
		TriggeringEvent:      triggeringEvent,
	}
	if lastPhone {
		message.RecipientPhoneNumber = user.LastPhone
	}
	return smsChannels.HandleMessage(message)
}
