package metrics

import (
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

func RegisterRequestCounter(recorder StatusRecorder, r *http.Request) {
	var labels = map[string]attribute.Value{
		URI:    attribute.StringValue(baseURI(r)),
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
		URI:        attribute.StringValue(baseURI(r)),
		Method:     attribute.StringValue(r.Method),
		ReturnCode: attribute.IntValue(recorder.Status()),
	}
	RegisterCounter(ReturnCodeCounter, ReturnCodeCounterDescription)
	AddCount(r.Context(), ReturnCodeCounter, 1, labels)
}

func baseURI(r *http.Request) string {
	return strings.Split(r.RequestURI, "?")[0]
}
