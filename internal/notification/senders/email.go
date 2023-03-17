package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
)

const smtpSpanName = "smtp.NotificationChannel"

func EmailChannels(
	ctx context.Context,
	emailConfig func(ctx context.Context) (*smtp.Config, error),
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	successMetricName,
	failureMetricName string,
) (chain *Chain, err error) {
	channels := make([]channels.NotificationChannel, 0, 3)
	p, err := smtp.InitChannel(ctx, emailConfig)
	// TODO: Why is this error not handled?
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
	channels = append(channels, debugChannels(ctx, getFileSystemProvider, getLogProvider)...)
	return chainChannels(channels...), nil
}
