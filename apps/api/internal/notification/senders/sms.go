package senders

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
)

const twilioSpanName = "twilio.NotificationChannel"

func SMSChannels(
	ctx context.Context,
	smsConfig *sms.Config,
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	successMetricName,
	failureMetricName string,
) (chain *Chain, err error) {
	channels := make([]channels.NotificationChannel, 0, 3)
	if smsConfig.TwilioConfig != nil {
		channels = append(
			channels,
			instrumenting.Wrap(
				ctx,
				twilio.InitChannel(*smsConfig.TwilioConfig),
				twilioSpanName,
				successMetricName,
				failureMetricName,
			),
		)
	}
	if smsConfig.WebhookConfig != nil {
		webhookChannel, err := webhook.InitChannel(ctx, *smsConfig.WebhookConfig)
		logging.WithFields(
			"instance", authz.GetInstance(ctx).InstanceID(),
			"callurl", smsConfig.WebhookConfig.CallURL,
		).OnError(err).Debug("initializing JSON channel failed")
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
	}
	channels = append(channels, debugChannels(ctx, getFileSystemProvider, getLogProvider)...)
	return ChainChannels(channels...), nil
}
