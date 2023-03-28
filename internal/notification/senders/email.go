package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
)

func EmailChannels(ctx context.Context, emailConfig func(ctx context.Context) (*smtp.Config, error), getFileSystemProvider func(ctx context.Context) (*fs.Config, error), getLogProvider func(ctx context.Context) (*log.Config, error)) (chain *Chain, err error) {
	channels := make([]channels.NotificationChannel, 0, 3)
	p, err := smtp.InitSMTPChannel(ctx, emailConfig)
	if err == nil {
		channels = append(channels, p)
	}
	channels = append(channels, debugChannels(ctx, getFileSystemProvider, getLogProvider)...)
	return chainChannels(channels...), nil
}
