package metrics

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
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
)

type Handler struct {
	handler http.Handler
	methods []MetricType
	filter  otelhttp.Filter
}

type MetricType int32

const (
	MetricTypeTotalCount MetricType = iota
	MetricTypeStatusCode
	MetricTypeRequestCount
)

type StatusRecorder struct {
	http.ResponseWriter
	RequestURI *string
	Status     int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type Filter func(*http.Request) bool

func NewHandler(handler http.Handler, metricMethods []MetricType, ignoredEndpoints ...string) http.Handler {
	h := Handler{
		handler: handler,
		methods: metricMethods,
		filter:  instrumentation.RequestFilter(ignoredEndpoints...),
	}
	return &h
}

type key int

const requestURI key = iota

func SetRequestURIPattern(ctx context.Context, pattern string) {
	uri, ok := ctx.Value(requestURI).(*string)
	if !ok {
		return
	}
	*uri = pattern
}

// ServeHTTP serves HTTP requests (http.Handler)
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
	uri := strings.Split(r.RequestURI, "?")[0]
	recorder := &StatusRecorder{
		ResponseWriter: w,
		RequestURI:     &uri,
		Status:         200,
	}
	r = r.WithContext(context.WithValue(r.Context(), requestURI, &uri))
	h.handler.ServeHTTP(recorder, r)
	if h.containsMetricsMethod(MetricTypeRequestCount) {
		RegisterRequestCounter(recorder, r)
	}
	if h.containsMetricsMethod(MetricTypeTotalCount) {
		RegisterTotalRequestCounter(r)
	}
	if h.containsMetricsMethod(MetricTypeStatusCode) {
		RegisterRequestCodeCounter(recorder, r)
	}
}

func (h *Handler) containsMetricsMethod(method MetricType) bool {
	for _, m := range h.methods {
		if m == method {
			return true
		}
	}
	return false
}

func RegisterRequestCounter(recorder *StatusRecorder, r *http.Request) {
	var labels = map[string]attribute.Value{
		URI:    attribute.StringValue(*recorder.RequestURI),
		Method: attribute.StringValue(r.Method),
	}
	RegisterCounter(RequestCounter, RequestCountDescription)
	AddCount(r.Context(), RequestCounter, 1, labels)
}

func RegisterTotalRequestCounter(r *http.Request) {
	RegisterCounter(TotalRequestCounter, TotalRequestDescription)
	AddCount(r.Context(), TotalRequestCounter, 1, nil)
}

func RegisterRequestCodeCounter(recorder *StatusRecorder, r *http.Request) {
	var labels = map[string]attribute.Value{
		URI:        attribute.StringValue(*recorder.RequestURI),
		Method:     attribute.StringValue(r.Method),
		ReturnCode: attribute.IntValue(recorder.Status),
	}
	RegisterCounter(ReturnCodeCounter, ReturnCodeCounterDescription)
	AddCount(r.Context(), ReturnCodeCounter, 1, labels)
}
