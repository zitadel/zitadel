package instrumenting

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
)

func logMessages[T channels.Message](ctx context.Context, channel channels.NotificationChannel[T]) channels.NotificationChannel[T] {
	return channels.HandleMessageFunc[T](func(message T) error {
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
