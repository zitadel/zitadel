package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
)

func SMSChannels(ctx context.Context, twilioConfig *twilio.TwilioConfig, getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error)) (chain *Chain, err error) {
	channels := make([]channels.NotificationChannel, 0, 3)
	if twilioConfig != nil {
		channels = append(channels, twilio.InitTwilioChannel(*twilioConfig))
	}
	channels = append(channels, debugChannels(ctx, getFileSystemProvider, getLogProvider)...)
	return chainChannels(channels...), nil
}
