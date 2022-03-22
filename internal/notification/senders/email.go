package senders

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/notification/channels/log"
	"github.com/caos/zitadel/internal/notification/channels/smtp"
)

func EmailChannels(ctx context.Context, config systemdefaults.Notifications, emailConfig func(ctx context.Context) (*smtp.EmailConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error)) (chain *Chain, err error) {
	p, err := smtp.InitSMTPChannel(ctx, emailConfig)
	if err == nil {
		chain.channels = append(chain.channels, p)
	}
	chain, err = debugChannels(ctx, config, getFileSystemProvider, getLogProvider)
	if err != nil {
		logging.New().Info("Error in creating debug channels")
	}
	return chain, nil
}
