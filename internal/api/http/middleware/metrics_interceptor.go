package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
	// Requested paths contain object IDs and are therefore unbounded, so the router below
	// gets the opportunity to report the route pattern it served the request with instead.
	ctx := metrics.WithRequestURIPattern(r.Context())
	r = r.WithContext(ctx)
	h.handler.ServeHTTP(recorder, r)
	if pattern := chiRoutePattern(r); pattern != "" {
		metrics.SetRequestURIPattern(ctx, pattern)
	}
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

// chiRoutePattern returns the route pattern chi matched the request against, e.g.
// "/oauth/v2/register/{client_id}", or [metrics.UnknownPath] if chi routed the request but
// found no route for it. Routers that are not chi based report their pattern themselves
// through [metrics.SetRequestURIPattern]; for those an empty string is returned.
//
// chi only knows the pattern once it routed the request, so this must be called after the
// wrapped handler returned.
func chiRoutePattern(r *http.Request) string {
	routeCtx := chi.RouteContext(r.Context())
	if routeCtx == nil {
		return ""
	}
	if pattern := routeCtx.RoutePattern(); pattern != "" {
		return pattern
	}
	return metrics.UnknownPath
}

func (h *Handler) containsMetricsMethod(method metrics.MetricType) bool {
	for _, m := range h.methods {
		if m == method {
			return true
		}
	}
	return false
}
