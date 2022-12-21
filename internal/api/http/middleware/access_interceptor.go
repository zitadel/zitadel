package middleware

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zitadel/zitadel/internal/logstore/emitters/access"

	"github.com/zitadel/zitadel/internal/logstore"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
)

const (
	// TODO: Make configurable?
	limitExceededCookiename   = "zitadel.quota.exceeded"
	limitExceededCookieMaxAge = 60 * 5 // 5 minutes
)

type AccessInterceptor struct {
	svc           *logstore.Service
	cookieHandler *http_utils.CookieHandler
}

func NewAccessInterceptor(svc *logstore.Service) *AccessInterceptor {
	return &AccessInterceptor{
		svc: svc,
		cookieHandler: http_utils.NewCookieHandler(
			http_utils.WithUnsecure(),
			http_utils.WithMaxAge(limitExceededCookieMaxAge),
		),
	}
}

func (a *AccessInterceptor) Handle(next http.Handler) http.Handler {
	if !a.svc.Enabled() {
		return next
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		wrappedWriter := &statusRecorder{ResponseWriter: writer, status: 0}

		ctx := request.Context()
		instance := authz.GetInstance(ctx)
		limit, err := a.svc.Limit(ctx, instance.InstanceID())
		if err != nil {
			logging.Warnf("failed to check whether requests should be limited: %s", err.Error())
			err = nil
		}

		a.cookieHandler.SetCookie(wrappedWriter, limitExceededCookiename, request.Host, strconv.FormatBool(limit))

		if limit {
			wrappedWriter.WriteHeader(http.StatusTooManyRequests)
			wrappedWriter.ignoreWrites = true
		}

		next.ServeHTTP(wrappedWriter, request)

		requestURL := request.RequestURI
		unescapedURL, err := url.QueryUnescape(requestURL)
		if err != nil {
			logging.Warningf("failed to unescape request url %s", requestURL)
		}
		err = a.svc.Handle(ctx, &access.Record{
			Timestamp:       time.Now(),
			Protocol:        access.HTTP,
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
	status       int
	ignoreWrites bool
}

func (r *statusRecorder) WriteHeader(status int) {
	if r.ignoreWrites {
		return
	}
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
