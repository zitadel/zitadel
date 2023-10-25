package senders

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/resources"
)

const webhookSpanName = "webhook.NotificationChannel"

func WebhookChannels(
	ctx context.Context,
	queries *resources.NotificationQueries,
	webhookConfig webhook.Config,
	successMetricName,
	failureMetricName string,
) (*Chain[*messages.JSON], error) {
	if err := webhookConfig.Validate(); err != nil {
		return nil, err
	}
	channels := make([]channels.NotificationChannel[*messages.JSON], 0, 3)
	webhookChannel, err := webhook.Connect(ctx, webhookConfig)
	logging.WithFields(
		"instance", authz.GetInstance(ctx).InstanceID(),
		"callurl", webhookConfig.CallURL,
	).OnError(err).Debug("connecting to JSON channel failed")
	if err == nil {
		channels = append(
			channels,
			instrumenting.Wrap(
				ctx,
				webhookChannel,
				webhookSpanName,
				successMetricName,
				failureMetricName,
			),
		)
	}
	channels = append(channels, connectToDebugChannels[*messages.JSON](ctx, queries)...)
	return ChainChannels(channels...), nil
}
