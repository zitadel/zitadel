package senders

import (
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/channels"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/notification/channels/log"
)

func debugChannels(config systemdefaults.Notifications) (channels.NotificationChannel, error) {
	var (
		providers []channels.NotificationChannel
	)

	if config.Providers.FileSystem.Enabled {
		p, err := fs.InitFSChannel(config.Providers.FileSystem)
		if err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}

	if config.Providers.Log.Enabled {
		providers = append(providers, log.InitStdoutChannel(config.Providers.Log))
	}

	return chainChannels(providers...), nil
}
