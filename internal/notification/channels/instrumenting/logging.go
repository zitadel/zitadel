package instrumenting

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/notification/channels"
)

func logMessages(ctx context.Context, channel channels.NotificationChannel) channels.NotificationChannel {
	return channels.HandleMessageFunc(func(message channels.Message) error {
		logEntry := logging.WithFields(
			"instance", authz.GetInstance(ctx).InstanceID(),
			"triggering_event_type", message.GetTriggeringEvent().Type(),
		)
		logEntry.Debug("sending notification")
		err := channel.HandleMessage(message)
		logEntry.OnError(err).Warn("sending notification failed")
		logEntry.Debug("notification sent")
		return err
	})
}
