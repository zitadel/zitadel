package middleware

import (
	"net/http"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"

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

		requestURL := request.RequestURI
		unescapedURL, err := url.QueryUnescape(requestURL)
		if err != nil {
			logging.Warningf("failed to unescape request url %s", requestURL)
		}

		ctx := request.Context()
		instance := authz.GetInstance(ctx)
		limit, err := a.svc.Limit(ctx, instance.InstanceID())
		if err != nil {
			logging.Warnf("failed to check whether requests should be limited: %s", err.Error())
			err = nil
		}

		if limit {
			wrappedWriter.WriteHeader(http.StatusTooManyRequests)
		}

		next.ServeHTTP(wrappedWriter, request)

		err = a.svc.Handle(ctx, &logstore.AccessLogRecord{
			Timestamp:       time.Now(),
			Protocol:        logstore.HTTP,
			RequestURL:      unescapedURL,
			ResponseStatus:  uint32(wrappedWriter.status),
			RequestHeaders:  request.Header,
			ResponseHeaders: writer.Header(),
			InstanceID:      instance.InstanceID(),
			ProjectID:       instance.ProjectID(),
			RequestedDomain: instance.RequestedDomain(),
			RequestedHost:   instance.RequestedHost(),
		})

		if err != nil {
			logging.Warnf("failed to handle access log: %s", err.Error())
			err = nil
		}
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
