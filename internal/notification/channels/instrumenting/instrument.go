package instrumenting

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/notification/channels"
)

func Wrap(
	ctx context.Context,
	channel channels.NotificationChannel,
	traceSpanName,
	successMetricName,
	failureMetricName string,
) channels.NotificationChannel {
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
