package instrumenting

import (
	"context"

	"github.com/zitadel/logging"
	"go.opentelemetry.io/otel/attribute"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/notification/channels"
	"github.com/zitadel/zitadel/v2/internal/telemetry/metrics"
)

func countMessages(ctx context.Context, channel channels.NotificationChannel, successMetricName, errorMetricName string) channels.NotificationChannel {
	return channels.HandleMessageFunc(func(message channels.Message) error {
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
		"triggering_event_typey": attribute.StringValue(string(message.GetTriggeringEvent().Type())),
		"instance":               attribute.StringValue(authz.GetInstance(ctx).InstanceID()),
	}
	if err != nil {
		labels["error"] = attribute.StringValue(err.Error())
	}
	addCountErr := metrics.AddCount(ctx, metricName, 1, labels)
	logging.WithFields("name", metricName, "labels", labels).OnError(addCountErr).Error("incrementing counter metric failed")
}
