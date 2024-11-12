package types

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	"github.com/zitadel/zitadel/internal/notification/messages"
)

func handleSecurityTokenEvent(
	ctx context.Context,
	setConfig set.Config,
	channels ChannelChains,
	token any,
	triggeringEvent eventstore.Event,
) error {
	message := &messages.Form{
		Serializable:    token,
		TriggeringEvent: triggeringEvent,
	}
	setChannels, err := channels.SecurityTokenEvent(ctx, setConfig)
	if err != nil {
		return err
	}
	return setChannels.HandleMessage(message)
}
