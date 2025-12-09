package instrumenting

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/tracing"
	"github.com/zitadel/zitadel/internal/notification/channels"
)

func traceMessages(ctx context.Context, channel channels.NotificationChannel, spanName string) channels.NotificationChannel {
	return channels.HandleMessageFunc(func(message channels.Message) (err error) {
		_, span := tracing.NewNamedSpan(ctx, spanName)
		defer func() { span.EndWithError(err) }()
		return channel.HandleMessage(message)
	})
}
