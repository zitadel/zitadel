package metrics

import (
	"net/http"
	"strings"

	"github.com/zitadel/logging"
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
)

type Handler struct {
	handler http.Handler
	methods []MetricType
	filters []Filter
}

type MetricType int32

const (
	MetricTypeTotalCount MetricType = iota
	MetricTypeStatusCode
	MetricTypeRequestCount
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type Filter func(*http.Request) bool

func NewMetricsHandler(handler http.Handler, metricMethods []MetricType, ignoredEndpoints ...string) http.Handler {
	h := Handler{
		handler: handler,
		methods: metricMethods,
	}
	if len(ignoredEndpoints) > 0 {
		h.filters = append(h.filters, shouldNotIgnore(ignoredEndpoints...))
	}
	return &h
}

// ServeHTTP serves HTTP requests (http.Handler)
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(h.methods) == 0 {
		h.handler.ServeHTTP(w, r)
		return
	}
	for _, f := range h.filters {
		if !f(r) {
			// Simply pass through to the handler if a filter rejects the request
			h.handler.ServeHTTP(w, r)
			return
		}
	}
	recorder := &StatusRecorder{
		ResponseWriter: w,
		Status:         200,
	}
	h.handler.ServeHTTP(recorder, r)
	if h.containsMetricsMethod(MetricTypeRequestCount) {
		RegisterRequestCounter(r)
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

func RegisterRequestCounter(r *http.Request) {
	var labels = map[string]attribute.Value{
		URI:    attribute.StringValue(strings.Split(r.RequestURI, "?")[0]),
		Method: attribute.StringValue(r.Method),
	}
	err := RegisterCounter(RequestCounter, RequestCountDescription)
	logging.WithFields("metric", RequestCounter).OnError(err).Panic("unable to register counter")
	err = AddCount(r.Context(), RequestCounter, 1, labels)
	logging.WithFields("metric", RequestCounter).OnError(err).Panic("unable to add count")
}

func RegisterTotalRequestCounter(r *http.Request) {
	err := RegisterCounter(TotalRequestCounter, TotalRequestDescription)
	logging.WithFields("metric", TotalRequestCounter).OnError(err).Panic("unable to register counter")
	err = AddCount(r.Context(), TotalRequestCounter, 1, nil)
	logging.WithFields("metric", TotalRequestCounter).OnError(err).Panic("unable to add count")
}

func RegisterRequestCodeCounter(recorder *StatusRecorder, r *http.Request) {
	var labels = map[string]attribute.Value{
		URI:        attribute.StringValue(strings.Split(r.RequestURI, "?")[0]),
		Method:     attribute.StringValue(r.Method),
		ReturnCode: attribute.IntValue(recorder.Status),
	}
	err := RegisterCounter(ReturnCodeCounter, ReturnCodeCounterDescription)
	logging.WithFields("metric", ReturnCodeCounter).OnError(err).Panic("unable to register counter")
	err = AddCount(r.Context(), ReturnCodeCounter, 1, labels)
	logging.WithFields("metric", ReturnCodeCounter).OnError(err).Panic("unable to add count")
}

func shouldNotIgnore(endpoints ...string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		for _, endpoint := range endpoints {
			if strings.HasPrefix(r.URL.RequestURI(), endpoint) {
				return false
			}
		}
		return true
	}
}
