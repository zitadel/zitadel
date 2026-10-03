package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
)

func DefaultTraceHandler(handler http.Handler) http.Handler {
	return TraceHandler(false, http_utils.Probes...)(handler)
}

// TraceHandler wraps handler with otelhttp instrumentation.
//
// When trustRemoteSpans is false (the default), every request is treated as
// a [otelhttp.WithPublicEndpointFn] endpoint: an incoming traceparent header
// is never used to continue the caller's trace, it is at most recorded as a
// span link. Set trustRemoteSpans to true for handlers where the caller is
// known and trusted (mirroring the behavior already available for the
// gRPC/Connect API via ZITADEL_INSTRUMENTATION_TRACE_TRUSTREMOTESPANS, see
// otelconnect.WithTrustRemote in internal/api/api.go) so that spans emitted
// by this handler are correctly parented under the caller's trace instead of
// starting a disconnected trace.
func TraceHandler(trustRemoteSpans bool, ignoredPrefix ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return otelhttp.NewHandler(handler,
			"zitadel",
			otelhttp.WithFilter(instrumentation.RequestFilter(ignoredPrefix...)),
			otelhttp.WithPublicEndpointFn(func(_ *http.Request) bool {
				return !trustRemoteSpans
			}),
			otelhttp.WithSpanNameFormatter(spanNameFormatter),
			otelhttp.WithMeterProvider(otel.GetMeterProvider()),
		)
	}
}

func spanNameFormatter(_ string, r *http.Request) string {
	return r.URL.Path
}
