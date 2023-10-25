package senders

import (
	"context"
	"github.com/zitadel/zitadel/internal/notification/channels/email_webhook"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/resources"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
)

const (
	smtpSpanName         = "smtp.NotificationChannel"
	emailWebhookSpanName = "emailwebhook.NotificationChannel"
)

func EmailChannels(
	ctx context.Context,
	queries *resources.NotificationQueries,
	smtpConfig *smtp.Config,
	webhookConfig *email_webhook.Config,
	successMetricSMTPName,
	failureMetricSMTPName,
	successMetricEmailWebhookName,
	failureMetricEmailWebhookName string,
) (chain *Chain[*messages.Email], err error) {
	channels := make([]channels.NotificationChannel[*messages.Email], 0, 4)
	if smtpConfig != nil {
		p, err := smtp.Connect(smtpConfig)
		logging.WithFields(
			"instance", authz.GetInstance(ctx).InstanceID(),
		).OnError(err).Debug("connecting to SMTP channel failed")
		if err == nil {
			channels = append(channels, instrumenting.Wrap[*messages.Email](ctx, p, smtpSpanName, successMetricSMTPName, failureMetricSMTPName))
		}
	}
	if webhookConfig != nil {
		p, err := email_webhook.Connect(ctx, webhookConfig)
		logging.WithFields(
			"instance", authz.GetInstance(ctx).InstanceID(),
			"callurl", webhookConfig.Webhook.CallURL,
		).OnError(err).Debug("connecting to email webhook channel failed")
		if err == nil {
			channels = append(channels, instrumenting.Wrap[*messages.Email](ctx, p, emailWebhookSpanName, successMetricEmailWebhookName, failureMetricEmailWebhookName))
		}
	}
	channels = append(channels, connectToDebugChannels[*messages.Email](ctx, queries)...)
	return ChainChannels(channels...), nil
}
