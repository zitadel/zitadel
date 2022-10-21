package middleware

import (
	"net/http"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/internal/logstore/access"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/logstore"
)

type AccessInterceptor struct {
	svc *access.Service
}

func NewAccessInterceptor(svc *access.Service) *AccessInterceptor {
	return &AccessInterceptor{svc: svc}
}

func (a *AccessInterceptor) Handler(next http.Handler) http.Handler {
	if !a.svc.Enabled() {
		return next
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		wrappedWriter := &statusRecorder{ResponseWriter: writer, status: 0}
		next.ServeHTTP(wrappedWriter, request)

		requestURL := request.RequestURI
		unescapedURL, err := url.QueryUnescape(requestURL)
		if err != nil {
			logging.Warningf("failed to unescape request url %s", requestURL)
		}

		a.svc.Handle(request.Context(), &logstore.AccessLogRecord{
			Timestamp:       time.Now(),
			Protocol:        logstore.HTTP,
			RequestURL:      unescapedURL,
			ResponseStatus:  uint32(wrappedWriter.status),
			RequestHeaders:  request.Header,
			ResponseHeaders: writer.Header(),
		})
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
