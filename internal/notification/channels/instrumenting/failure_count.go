package instrumenting

import (
	"context"

	"github.com/zitadel/logging"
	"go.opentelemetry.io/otel/attribute"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

func CountFailure(ctx context.Context, counterMetricName string) {
	instanceID := authz.GetInstance(ctx).InstanceID()
	err := metrics.AddCount(
		ctx,
		counterMetricName,
		1,
		map[string]attribute.Value{
			"instance": attribute.StringValue(instanceID),
		},
	)
	logging.OnError(err).
		WithField("metring", counterMetricName).
		WithField("instance", instanceID).
		Error("incrementing counter failed")
}
