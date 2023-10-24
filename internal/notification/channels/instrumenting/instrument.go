package instrumenting

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
)

func Wrap[T channels.Message](
	ctx context.Context,
	channel channels.NotificationChannel[T],
	traceSpanName,
	successMetricName,
	failureMetricName string,
) channels.NotificationChannel[T] {
	return traceMessages(
		ctx,
		countMessages(
			ctx,
			logMessages(ctx, channel),
			successMetricName,
			failureMetricName,
		),
		traceSpanName,
	)
}
