package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
)

func MetricsHandler(metricTypes []metrics.MetricType, ignoredMethods ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return metrics.NewHandler(handler, metricTypes, ignoredMethods...)
	}
}
