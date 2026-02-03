package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
)

func MetricsHandler(metricTypes []metrics.MetricType, ignoredMethods ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &Handler{
			handler: handler,
			methods: metricTypes,
			filter:  instrumentation.RequestFilter(ignoredMethods...),
		}
	}
}

type Handler struct {
	handler http.Handler
	methods []metrics.MetricType
	filter  otelhttp.Filter
}

// ServeHTTP implements [http.Handler]
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(h.methods) == 0 {
		h.handler.ServeHTTP(w, r)
		return
	}
	if !h.filter(r) {
		// Simply pass through to the handler if a filter rejects the request
		h.handler.ServeHTTP(w, r)
		return
	}
	recorder := newStatusWriter(w)
	h.handler.ServeHTTP(recorder, r)
	if h.containsMetricsMethod(metrics.MetricTypeRequestCount) {
		metrics.RegisterRequestCounter(recorder, r)
	}
	if h.containsMetricsMethod(metrics.MetricTypeTotalCount) {
		metrics.RegisterTotalRequestCounter(r)
	}
	if h.containsMetricsMethod(metrics.MetricTypeStatusCode) {
		metrics.RegisterRequestCodeCounter(recorder, r)
	}
}

func (h *Handler) containsMetricsMethod(method metrics.MetricType) bool {
	for _, m := range h.methods {
		if m == method {
			return true
		}
	}
	return false
}
