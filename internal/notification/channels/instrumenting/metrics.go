package instrumenting

import (
	"context"

	"github.com/zitadel/logging"
	"go.opentelemetry.io/otel/attribute"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

func countMessages[T channels.Message](ctx context.Context, channel channels.NotificationChannel[T], successMetricName, errorMetricName string) channels.NotificationChannel[T] {
	return channels.HandleMessageFunc[T](func(message T) error {
		err := channel.HandleMessage(message)
		metricName := successMetricName
		if err != nil {
			metricName = errorMetricName
		}
		addCount(ctx, metricName, message, err)
		return err
	})
}

func addCount(ctx context.Context, metricName string, message channels.Message, err error) {
	labels := map[string]attribute.Value{
		"triggering_event_type": attribute.StringValue(string(message.GetTriggeringEvent().Type())),
		"instance":              attribute.StringValue(authz.GetInstance(ctx).InstanceID()),
	}
	if err != nil {
		labels["error"] = attribute.StringValue(err.Error())
	}
	addCountErr := metrics.AddCount(ctx, metricName, 1, labels)
	logging.WithFields("name", metricName, "labels", labels).OnError(addCountErr).Error("incrementing counter metric failed")
}
