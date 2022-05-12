package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
)

func SMSChannels(ctx context.Context, twilioConfig *twilio.TwilioConfig, getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error)) (chain *Chain, err error) {
	chain = debugChannels(ctx, getFileSystemProvider, getLogProvider)
	if twilioConfig != nil {
		chain.channels = append(chain.channels, twilio.InitTwilioChannel(*twilioConfig))
	}
	return chain, nil
}
