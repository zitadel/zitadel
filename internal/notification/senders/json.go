package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
)

func JSONChannels(
	ctx context.Context,
	webhookConfig webhook.Config,
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
) (*Chain, error) {
	if err := webhookConfig.IsValid(); err != nil {
		return nil, err
	}
	webhookChannel, err := webhook.InitWebhookChannel(ctx, webhookConfig)
	if err != nil {
		return nil, err
	}
	channels := []channels.NotificationChannel{webhookChannel}
	channels = append(channels, debugChannels(ctx, getFileSystemProvider, getLogProvider)...)
	return chainChannels(channels...), nil
}
