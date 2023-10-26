package senders

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
)

const smtpSpanName = "smtp.NotificationChannel"

func EmailChannels(
	ctx context.Context,
	emailConfig *smtp.Config,
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	successMetricName,
	failureMetricName string,
) (chain *Chain, err error) {
	channels := make([]channels.NotificationChannel, 0, 3)
	p, err := smtp.InitChannel(emailConfig)
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
	channels = append(channels, debugChannels(ctx, getFileSystemProvider, getLogProvider)...)
	return ChainChannels(channels...), nil
}
