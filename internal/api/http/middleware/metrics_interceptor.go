package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/v2/internal/telemetry/metrics"
)

func MetricsHandler(metricTypes []metrics.MetricType, ignoredMethods ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return metrics.NewMetricsHandler(handler, metricTypes, ignoredMethods...)
	}
}
