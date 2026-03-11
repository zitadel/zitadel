package signals

import (
	"context"
	"net/http"
	"strings"
	"time"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

// signalServicePrefix is the gRPC package prefix for the signal API itself.
// Calls to this service are excluded from request-stream signals to avoid
// the signal explorer polluting its own data.
const signalServicePrefix = "/zitadel.signal."

// SignalConnectUnaryInterceptor returns a ConnectRPC unary interceptor that
// emits a fire-and-forget risk signal after every call. If the emitter is nil
// the interceptor is a no-op pass-through.
func SignalConnectUnaryInterceptor(emitter *Emitter, geoCountryHeader string) connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			resp, handlerErr := handler(ctx, req)

			if emitter == nil {
				return resp, handlerErr
			}

			// Skip self-recording: do not emit signals for the signal API itself.
			if strings.HasPrefix(req.Spec().Procedure, signalServicePrefix) {
				return resp, handlerErr
			}

			ctxData := authz.GetCtxData(ctx)
			instance := authz.GetInstance(ctx)

			outcome := OutcomeSuccess
			if handlerErr != nil {
				outcome = OutcomeFailure
			}

			// Convert ConnectRPC headers to http.Header for extraction.
			hctx := ExtractHTTPContext(http.Header(req.Header()), geoCountryHeader)

			sig := Signal{
				InstanceID:     instance.InstanceID(),
				UserID:         ctxData.UserID,
				CallerID:       ctxData.UserID,
				Stream:         StreamRequests,
				Operation:      req.Spec().Procedure,
				IP:             http_util.RemoteIPFromCtx(ctx),
				UserAgent:      truncateString(req.Header().Get(http_util.UserAgentHeader), maxUserAgentLen),
				Outcome:        outcome,
				Timestamp:      time.Now().UTC(),
				AcceptLanguage: hctx.AcceptLanguage,
				Country:        hctx.Country,
				ForwardedChain: hctx.ForwardedChain,
				Referer:        hctx.Referer,
				SecFetchSite:   hctx.SecFetchSite,
				IsHTTPS:        hctx.IsHTTPS,
				TraceID:        tracing.TraceIDFromCtx(ctx),
				SpanID:         tracing.SpanIDFromCtx(ctx),
			}
			emitter.Emit(sig)
			return resp, handlerErr
		}
	}
}

// SignalHTTPMiddleware returns an HTTP middleware that emits a fire-and-forget
// signal for every request. It covers all paths not served by connectRPC
// (OIDC, SAML, login UI, health checks, etc.).
//
// If emitter is nil the middleware is a transparent pass-through.
// The signal is emitted *after* the downstream handler returns so the
// outcome (HTTP status code) is available.
func SignalHTTPMiddleware(emitter *Emitter, geoCountryHeader string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if emitter == nil {
				next.ServeHTTP(w, r)
				return
			}

			rw := &statusCapture{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			// Skip self-recording: do not emit signals for the signal API routes.
			if strings.HasPrefix(r.URL.Path, "/zitadel.signal.") || strings.HasPrefix(r.URL.Path, "/api/v2/signals") {
				return
			}

			ctx := r.Context()
			instance := authz.GetInstance(ctx)
			ctxData := authz.GetCtxData(ctx)

			outcome := OutcomeSuccess
			if rw.status >= 400 {
				outcome = OutcomeFailure
			}

			hctx := ExtractHTTPContext(r.Header, geoCountryHeader)
			emitter.Emit(Signal{
				InstanceID:     instance.InstanceID(),
				UserID:         ctxData.UserID,
				CallerID:       ctxData.UserID,
				Stream:         StreamRequests,
				Operation:      r.Method + " " + r.URL.Path,
				IP:             http_util.RemoteIPFromCtx(ctx),
				UserAgent:      truncateString(r.Header.Get("User-Agent"), maxUserAgentLen),
				Outcome:        outcome,
				Timestamp:      time.Now().UTC(),
				AcceptLanguage: hctx.AcceptLanguage,
				Country:        hctx.Country,
				ForwardedChain: hctx.ForwardedChain,
				Referer:        hctx.Referer,
				SecFetchSite:   hctx.SecFetchSite,
				IsHTTPS:        hctx.IsHTTPS,
				TraceID:        tracing.TraceIDFromCtx(ctx),
				SpanID:         tracing.SpanIDFromCtx(ctx),
			})
		})
	}
}

// statusCapture wraps http.ResponseWriter to capture the written status code.
type statusCapture struct {
	http.ResponseWriter
	status int
}

func (s *statusCapture) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
