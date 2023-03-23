package types

import (
	"context"

	"github.com/zitadel/logging"

	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/query"
)

func generateSms(ctx context.Context, user *query.NotifyUser, content string, getTwilioProvider func(ctx context.Context) (*twilio.Config, error), getFileSystemProvider func(ctx context.Context) (*fs.Config, error), getLogProvider func(ctx context.Context) (*log.Config, error), lastPhone bool) error {
	number := ""
	twilioConfig, err := getTwilioProvider(ctx)
	if err == nil {
		number = twilioConfig.SenderNumber
	}
	message := &messages.SMS{
		SenderPhoneNumber:    number,
		RecipientPhoneNumber: user.VerifiedPhone,
		Content:              content,
	}
	if lastPhone {
		message.RecipientPhoneNumber = user.LastPhone
	}

	channelChain, err := senders.SMSChannels(ctx, twilioConfig, getFileSystemProvider, getLogProvider)
	logging.OnError(err).Error("could not create sms channel")

	if channelChain.Len() == 0 {
		return caos_errors.ThrowPreconditionFailed(nil, "PHONE-w8nfow", "Errors.Notification.Channels.NotPresent")
	}
	return channelChain.HandleMessage(message)
}
