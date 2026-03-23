package signals

// PREVIEW: Identity Signals is a preview feature. APIs, storage format,
// and configuration may change between releases without notice.

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/felixge/httpsnoop"
	otel_trace "go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

// signalServicePrefix is the gRPC package prefix for the signal API itself.
// Calls to this service are excluded to avoid self-recording.
const signalServicePrefix = "/zitadel.signal."

// SignalConnectUnaryInterceptor returns a ConnectRPC unary interceptor that
// emits a fire-and-forget signal after every call. If the emitter is nil
// the interceptor is a no-op pass-through.
func SignalConnectUnaryInterceptor(emitter *Emitter, geoCountryHeader string) connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			resp, handlerErr := handler(ctx, req)

			if emitter == nil {
				return resp, handlerErr
			}

			if strings.HasPrefix(req.Spec().Procedure, signalServicePrefix) {
				return resp, handlerErr
			}

			ctxData := authz.GetCtxData(ctx)
			instance := authz.GetInstance(ctx)

			outcome := OutcomeSuccess
			if handlerErr != nil {
				outcome = OutcomeFailure
			}

			hctx := ExtractHTTPContext(http.Header(req.Header()), geoCountryHeader)

			emitter.Emit(Signal{
				InstanceID:     instance.InstanceID(),
				UserID:         ctxData.UserID,
				CallerID:       ctxData.UserID,
				OrgID:          ctxData.OrgID,
				ProjectID:      ctxData.ProjectID,
				ClientID:       ctxData.AgentID,
				Stream:         StreamRequests,
				Operation:      req.Spec().Procedure,
				IP:             stripPort(http_util.RemoteIPFromCtx(ctx)),
				UserAgent:      truncateString(req.Header().Get(http_util.UserAgentHeader), maxUserAgentLen),
				Outcome:        outcome,
				Timestamp:      start.UTC(),
				DurationMs:     time.Since(start).Milliseconds(),
				AcceptLanguage: hctx.AcceptLanguage,
				Country:        hctx.Country,
				ForwardedChain: hctx.ForwardedChain,
				Referer:        hctx.Referer,
				SecFetchSite:   hctx.SecFetchSite,
				IsHTTPS:        hctx.IsHTTPS,
				TraceID:        tracing.TraceIDFromCtx(ctx),
				SpanID:         spanIDFromCtx(ctx),
			})
			return resp, handlerErr
		}
	}
}

// SignalHTTPMiddleware returns an HTTP middleware that emits a fire-and-forget
// signal for every request. Covers OIDC, SAML, login UI, health checks, etc.
//
// If emitter is nil the middleware is a transparent pass-through.
func SignalHTTPMiddleware(emitter *Emitter, geoCountryHeader string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if emitter == nil {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			metrics := httpsnoop.CaptureMetrics(next, w, r)

			if strings.HasPrefix(r.URL.Path, "/zitadel.signal.") ||
				strings.HasPrefix(r.URL.Path, "/v2/signals") ||
				strings.HasPrefix(r.URL.Path, "/api/v2/signals") {
				return
			}

			// Skip ConnectRPC and gRPC requests — covered by SignalConnectUnaryInterceptor.
			ct := r.Header.Get("Content-Type")
			if strings.HasPrefix(ct, "application/connect+") ||
				strings.HasPrefix(ct, "application/grpc") {
				return
			}

			ctx := r.Context()
			instance := authz.GetInstance(ctx)
			ctxData := authz.GetCtxData(ctx)

			outcome := OutcomeSuccess
			if metrics.Code >= 400 {
				outcome = OutcomeFailure
			}

			hctx := ExtractHTTPContext(r.Header, geoCountryHeader)
			emitter.Emit(Signal{
				InstanceID:     instance.InstanceID(),
				UserID:         ctxData.UserID,
				CallerID:       ctxData.UserID,
				OrgID:          ctxData.OrgID,
				ProjectID:      ctxData.ProjectID,
				ClientID:       ctxData.AgentID,
				Stream:         StreamRequests,
				Operation:      r.Method + " " + r.URL.Path,
				IP:             stripPort(http_util.RemoteIPFromCtx(ctx)),
				UserAgent:      truncateString(r.Header.Get("User-Agent"), maxUserAgentLen),
				Outcome:        outcome,
				Timestamp:      start.UTC(),
				DurationMs:     metrics.Duration.Milliseconds(),
				AcceptLanguage: hctx.AcceptLanguage,
				Country:        hctx.Country,
				ForwardedChain: hctx.ForwardedChain,
				Referer:        hctx.Referer,
				SecFetchSite:   hctx.SecFetchSite,
				IsHTTPS:        hctx.IsHTTPS,
				TraceID:        tracing.TraceIDFromCtx(ctx),
				SpanID:         spanIDFromCtx(ctx),
			})
		})
	}
}

// spanIDFromCtx extracts the OpenTelemetry span ID from the context.
func spanIDFromCtx(ctx context.Context) string {
	sc := otel_trace.SpanFromContext(ctx).SpanContext()
	if sc.HasSpanID() {
		return sc.SpanID().String()
	}
	return ""
}

// stripPort returns the IP address without the port suffix.
func stripPort(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}

