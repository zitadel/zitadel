package types

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"
)

func handleJSON(
	ctx context.Context,
	webhookConfig webhook.Config,
	getFileSystemProvider func(ctx context.Context) (*fs.Config, error),
	getLogProvider func(ctx context.Context) (*log.Config, error),
	serializable interface{},
	triggeringEvent eventstore.Event,
	successMetricName,
	failureMetricName string,
) error {
	message := &messages.JSON{
		Serializable:    serializable,
		TriggeringEvent: triggeringEvent,
	}
	channelChain, err := senders.JSONChannels(
		ctx,
		webhookConfig,
		getFileSystemProvider,
		getLogProvider,
		successMetricName,
		failureMetricName,
	)
	if err != nil {
		return err
	}
	return channelChain.HandleMessage(message)
}
