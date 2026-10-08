package metrics

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
)

const (
	RequestCounter               = "http.server.request_count"
	RequestCountDescription      = "Request counter"
	TotalRequestCounter          = "http.server.total_request_count"
	TotalRequestDescription      = "Total return code counter"
	ReturnCodeCounter            = "http.server.return_code_counter"
	ReturnCodeCounterDescription = "Return code counter"
	Method                       = "method"
	URI                          = "uri"
	ReturnCode                   = "return_code"

	// UnknownPath is recorded as uri for requests that did not match any route.
	// Recording the requested path instead would allow anyone to create an unbounded
	// number of metric series by calling arbitrary paths.
	UnknownPath = "UNKNOWN_PATH"
)

type MetricType int32

const (
	MetricTypeTotalCount MetricType = iota
	MetricTypeStatusCode
	MetricTypeRequestCount
)

type StatusRecorder interface {
	Status() int
}

type requestURIKey struct{}

// WithRequestURIPattern returns a context that carries the uri label of the HTTP metrics
// recorded for the current request. The metrics middleware must call this before passing
// the request on, so that the router serving it can substitute its route pattern through
// [SetRequestURIPattern].
func WithRequestURIPattern(ctx context.Context) context.Context {
	var pattern string
	return context.WithValue(ctx, requestURIKey{}, &pattern)
}

// SetRequestURIPattern records pattern as the uri label of the current request's HTTP
// metrics, e.g. "/v2/sessions/{session_id}" instead of "/v2/sessions/385063742120926058".
//
// Routers that serve templated paths must call this for every response they produce,
// including errors. Without it every resource ID gets its own metric series and the
// number of series only grows with the number of objects in the installation.
func SetRequestURIPattern(ctx context.Context, pattern string) {
	uri, ok := ctx.Value(requestURIKey{}).(*string)
	if !ok {
		return
	}
	*uri = pattern
}

// RequestURI returns the uri label to record for r: the route pattern reported by the
// router through [SetRequestURIPattern], or the requested path if the router did not
// report one.
func RequestURI(r *http.Request) string {
	if pattern, ok := r.Context().Value(requestURIKey{}).(*string); ok && *pattern != "" {
		return *pattern
	}
	return strings.Split(r.RequestURI, "?")[0]
}

func RegisterRequestCounter(recorder StatusRecorder, r *http.Request) {
	var labels = map[string]attribute.Value{
		URI:    attribute.StringValue(RequestURI(r)),
		Method: attribute.StringValue(r.Method),
	}
	RegisterCounter(RequestCounter, RequestCountDescription)
	AddCount(r.Context(), RequestCounter, 1, labels)
}

func RegisterTotalRequestCounter(r *http.Request) {
	RegisterCounter(TotalRequestCounter, TotalRequestDescription)
	AddCount(r.Context(), TotalRequestCounter, 1, nil)
}

func RegisterRequestCodeCounter(recorder StatusRecorder, r *http.Request) {
	var labels = map[string]attribute.Value{
		URI:        attribute.StringValue(RequestURI(r)),
		Method:     attribute.StringValue(r.Method),
		ReturnCode: attribute.IntValue(recorder.Status()),
	}
	RegisterCounter(ReturnCodeCounter, ReturnCodeCounterDescription)
	AddCount(r.Context(), ReturnCodeCounter, 1, labels)
}
