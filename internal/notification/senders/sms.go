package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
)

const twilioSpanName = "twilio.NotificationChannel"

func SMSChannels(
	ctx context.Context,
	twilioConfig *twilio.Config,
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	successMetricName,
	failureMetricName string,
) (chain *Chain, err error) {
	channels := make([]channels.NotificationChannel, 0, 3)
	if twilioConfig != nil {
		channels = append(
			channels,
			instrumenting.Wrap(
				ctx,
				twilio.InitChannel(*twilioConfig),
				twilioSpanName,
				successMetricName,
				failureMetricName,
			),
		)
	}
	channels = append(channels, debugChannels(ctx, getFileSystemProvider, getLogProvider)...)
	return ChainChannels(channels...), nil
}
