package metrics

import (
	"net/http"
	"strings"
)

const (
	RequestCount                 = "http.server.request_count"
	RequestCountDescription      = "Request counter"
	ReturnCodeCounter            = "http.server.return_code_counter"
	ReturnCodeCounterDescription = "Return code counter"
	Method                       = "method"
	URI                          = "uri"
)

type Handler struct {
	handler http.Handler
	filters []Filter
}

type Filter func(*http.Request) bool

func NewMetricsHandler(handler http.Handler, ignoredEndpoints ...string) http.Handler {
	h := Handler{
		handler: handler,
	}
	h.filters = append(h.filters, shouldNotIgnore(ignoredEndpoints...))
	return &h
}

// ServeHTTP serves HTTP requests (http.Handler)
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, f := range h.filters {
		if !f(r) {
			// Simply pass through to the handler if a filter rejects the request
			h.handler.ServeHTTP(w, r)
			return
		}
	}
	registerRequestCounter(r)
}

func registerRequestCounter(r *http.Request) {
	var labels = map[string]interface{}{
		URI:    r.RequestURI,
		Method: r.Method,
	}
	M.RegisterCounter(RequestCount, RequestCountDescription)
	M.AddCount(r.Context(), RequestCount, 1, labels)
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
