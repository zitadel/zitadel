package senders

import (
	"context"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/instrumenting"
)

func instrument(
	ctx context.Context,
	channel channels.NotificationChannel,
	successMetricName,
	failureMetricName string,
) channels.NotificationChannel {
	return instrumenting.CountReturnValues(
		ctx,
		instrumenting.LogReturnValues(ctx, channel),
		successMetricName,
		failureMetricName,
	)
}
