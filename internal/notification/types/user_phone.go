package types

import (
	"context"

	caos_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/notification/channels/log"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/notification/messages"
	"github.com/caos/zitadel/internal/notification/senders"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

func generateSms(ctx context.Context, user *view_model.NotifyUser, content string, getTwilioProvider func(ctx context.Context) (*twilio.TwilioConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error), lastPhone bool) error {
	number := ""
	twilio, err := getTwilioProvider(ctx)
	if err == nil {
		number = twilio.SenderNumber
	}
	message := &messages.SMS{
		SenderPhoneNumber:    number,
		RecipientPhoneNumber: user.VerifiedPhone,
		Content:              content,
	}
	if lastPhone {
		message.RecipientPhoneNumber = user.LastPhone
	}

	channelChain, err := senders.SMSChannels(ctx, twilio, getFileSystemProvider, getLogProvider)

	if channelChain.Len() == 0 {
		return caos_errors.ThrowPreconditionFailed(nil, "PHONE-w8nfow", "Errors.Notification.Channels.NotPresent")
	}
	return channelChain.HandleMessage(message)
}
