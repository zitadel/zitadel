package senders

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
)

const smtpSpanName = "smtp.NotificationChannel"

func EmailChannels(
	ctx context.Context,
	emailConfig *email.Config,
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	successMetricName,
	failureMetricName string,
) (chain *Chain, err error) {
	channels := make([]channels.NotificationChannel, 0, 3)
	if emailConfig.SMTPConfig != nil {
		p, err := smtp.InitChannel(emailConfig.SMTPConfig)
		logging.WithFields(
			"instance", authz.GetInstance(ctx).InstanceID(),
		).OnError(err).Debug("initializing SMTP channel failed")
		if err == nil {
			channels = append(
				channels,
				instrumenting.Wrap(
					ctx,
					p,
					smtpSpanName,
					successMetricName,
					failureMetricName,
				),
			)
		}
	}
	if emailConfig.WebhookConfig != nil {
		webhookChannel, err := webhook.InitChannel(ctx, *emailConfig.WebhookConfig)
		logging.WithFields(
			"instance", authz.GetInstance(ctx).InstanceID(),
			"callurl", emailConfig.WebhookConfig.CallURL,
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
