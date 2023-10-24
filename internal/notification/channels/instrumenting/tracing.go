package instrumenting

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func traceMessages[T channels.Message](ctx context.Context, channel channels.NotificationChannel[T], spanName string) channels.NotificationChannel[T] {
	return channels.HandleMessageFunc[T](func(message T) (err error) {
		_, span := tracing.NewNamedSpan(ctx, spanName)
		defer func() { span.EndWithError(err) }()
		return channel.HandleMessage(message)
	})
}
