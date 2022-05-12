package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
)

func EmailChannels(ctx context.Context, emailConfig func(ctx context.Context) (*smtp.EmailConfig, error), getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error)) (chain *Chain, err error) {
	chain = debugChannels(ctx, getFileSystemProvider, getLogProvider)
	p, err := smtp.InitSMTPChannel(ctx, emailConfig)
	if err == nil {
		chain.channels = append(chain.channels, p)
	}
	return chain, nil
}
