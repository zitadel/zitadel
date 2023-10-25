package senders

import (
	"context"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/resources"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
)

const twilioSpanName = "twilio.NotificationChannel"

func SMSChannels(
	ctx context.Context,
	queries *resources.NotificationQueries,
	twilioConfig *twilio.Config,
	successMetricName,
	failureMetricName string,
) (chain *Chain[*messages.SMS], err error) {
	channels := make([]channels.NotificationChannel[*messages.SMS], 0, 3)
	if twilioConfig != nil {
		channels = append(
			channels,
			instrumenting.Wrap(
				ctx,
				twilio.Connect(*twilioConfig),
				twilioSpanName,
				successMetricName,
				failureMetricName,
			),
		)
	}
	channels = append(channels, connectToDebugChannels[*messages.SMS](ctx, queries)...)
	return ChainChannels(channels...), nil
}
