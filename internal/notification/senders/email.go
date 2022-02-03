package senders

import (
	"context"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/channels"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
)

func EmailChannels(ctx context.Context, config systemdefaults.Notifications, emailConfig func(ctx context.Context) (*smtp.EmailConfig, error)) (channels.NotificationChannel, error) {

	debug, err := debugChannels(config)
	if err != nil {
		return nil, err
	}

	if !config.DebugMode {
		p, err := smtp.InitSMTPChannel(ctx, emailConfig)
		if err != nil {
			return nil, err
		}
		return chainChannels(debug, p), nil
	}

	return debug, nil
}
