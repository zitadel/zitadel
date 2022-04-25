package senders

import (
	"context"

	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/notification/channels/log"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
)

func SMSChannels(ctx context.Context, twilioConfig *twilio.TwilioConfig, getFileSystemProvider func(ctx context.Context) (*fs.FSConfig, error), getLogProvider func(ctx context.Context) (*log.LogConfig, error)) (chain *Chain, err error) {
	if twilioConfig != nil {
		chain.channels = append(chain.channels, twilio.InitTwilioChannel(*twilioConfig))
	}
	chain, err = debugChannels(ctx, getFileSystemProvider, getLogProvider)
	if err != nil {
		return nil, err
	}
	return chain, nil
}
