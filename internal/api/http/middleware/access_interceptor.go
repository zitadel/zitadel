package middleware

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/access"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type AccessInterceptor struct {
	svc           *logstore.Service
	cookieHandler *http_utils.CookieHandler
	limitConfig   *AccessConfig
}

type AccessConfig struct {
	ExhaustedCookieKey    string
	ExhaustedCookieMaxAge time.Duration
}

func NewAccessInterceptor(svc *logstore.Service, cookieConfig *AccessConfig) *AccessInterceptor {
	return &AccessInterceptor{
		svc: svc,
		cookieHandler: http_utils.NewCookieHandler(
			http_utils.WithUnsecure(),
			http_utils.WithMaxAge(int(math.Floor(cookieConfig.ExhaustedCookieMaxAge.Seconds()))),
		),
		limitConfig: cookieConfig,
	}
}

func (a *AccessInterceptor) Handle(next http.Handler) http.Handler {
	if !a.svc.Enabled() {
		return next
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		ctx := request.Context()
		var err error

		tracingCtx, span := tracing.NewServerInterceptorSpan(ctx)
		defer func() { span.EndWithError(err) }()

		wrappedWriter := &statusRecorder{ResponseWriter: writer, status: 0}

		instance := authz.GetInstance(ctx)
		remaining := a.svc.Limit(tracingCtx, instance.InstanceID())
		limit := remaining != nil && *remaining == 0

		a.cookieHandler.SetCookie(wrappedWriter, a.limitConfig.ExhaustedCookieKey, request.Host, strconv.FormatBool(limit))

		if limit {
			wrappedWriter.WriteHeader(http.StatusTooManyRequests)
			wrappedWriter.ignoreWrites = true
		}

		next.ServeHTTP(wrappedWriter, request)

		requestURL := request.RequestURI
		unescapedURL, err := url.QueryUnescape(requestURL)
		if err != nil {
			logging.WithError(err).WithField("url", requestURL).Warning("failed to unescape request url")
			// err = nil is effective because of deferred tracing span end
			err = nil
		}
		a.svc.Handle(tracingCtx, &access.Record{
			LogDate:         time.Now(),
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
