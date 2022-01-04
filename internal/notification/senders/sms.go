package senders

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/channels"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
)

func SMSChannels(config systemdefaults.Notifications) (channels.NotificationChannel, error) {

	debug, err := debugChannels(config)
	if err != nil {
		return nil, err
	}

	if !config.DebugMode {
		return chainChannels(twilio.InitTwilioProvider(config.Providers.Twilio), debug), nil
	}

	return debug, nil
}
