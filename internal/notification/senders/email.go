package senders

import (
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
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
		return chainChannels(debug, p), nil
	}

	return debug, nil
}
