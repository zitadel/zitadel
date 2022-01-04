package senders

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/channels"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
)

func EmailChannels(config systemdefaults.Notifications) (channels.NotificationChannel, error) {

	debug, err := debugChannels(config)
	if err != nil {
		return nil, err
	}

	if !config.DebugMode {
		p, err := smtp.InitSMTPChannel(config.Providers.Email)
		if err != nil {
			return nil, err
		}
		return chainChannels(p, debug), nil
	}

	return debug, nil
}
