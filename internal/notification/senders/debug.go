package senders

import (
	"context"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/channels"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/notification/channels/log"
)

func debugChannels(ctx context.Context, config systemdefaults.Notifications, getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error)) (*Chain, error) {
	var (
		providers []channels.NotificationChannel
	)

	if fsProvider, err := getFileSystemProvider(ctx); err == nil {
		p, err := fs.InitFSChannel(config.FileSystemPath, *fsProvider)
		if err == nil {
			providers = append(providers, p)
		}
	}

	if logProvider, err := getLogProvider(ctx); err == nil {
		providers = append(providers, log.InitStdoutChannel(*logProvider))
	}

	return chainChannels(providers...), nil
}
