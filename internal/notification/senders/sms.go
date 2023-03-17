package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
)

func SMSChannels(ctx context.Context, twilioConfig *twilio.Config, getFileSystemProvider func(ctx context.Context) (*fs.Config, error), getLogProvider func(ctx context.Context) (*log.Config, error)) (chain *Chain, err error) {
	channels := make([]channels.NotificationChannel, 0, 3)
	if twilioConfig != nil {
		channels = append(channels, twilio.InitTwilioChannel(*twilioConfig))
	}
	channels = append(channels, debugChannels(ctx, getFileSystemProvider, getLogProvider)...)
	return chainChannels(channels...), nil
}
