package middleware

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/server/middleware"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/emitters/access"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type AccessInterceptor struct {
	svc           *logstore.Service
	cookieHandler *http_utils.CookieHandler
	limitConfig   *AccessConfig
	storeOnly     bool
}

type AccessConfig struct {
	ExhaustedCookieKey    string
	ExhaustedCookieMaxAge time.Duration
}

// NewAccessInterceptor intercepts all requests and stores them to the logstore.
// If storeOnly is false, it also checks if requests are exhausted.
// If requests are exhausted, it also returns http.StatusTooManyRequests and sets a cookie
func NewAccessInterceptor(svc *logstore.Service, cookieHandler *http_utils.CookieHandler, cookieConfig *AccessConfig) *AccessInterceptor {
	return &AccessInterceptor{
		svc:           svc,
		cookieHandler: cookieHandler,
		limitConfig:   cookieConfig,
	}
}

func (a *AccessInterceptor) WithoutLimiting() *AccessInterceptor {
	return &AccessInterceptor{
		svc:           a.svc,
		cookieHandler: a.cookieHandler,
		limitConfig:   a.limitConfig,
		storeOnly:     true,
	}
}

func (a *AccessInterceptor) AccessService() *logstore.Service {
	return a.svc
}

func (a *AccessInterceptor) Limit(ctx context.Context) bool {
	if !a.svc.Enabled() || a.storeOnly {
		return false
	}
	instance := authz.GetInstance(ctx)
	remaining := a.svc.Limit(ctx, instance.InstanceID())
	return remaining != nil && *remaining <= 0
}

func (a *AccessInterceptor) SetExhaustedCookie(writer http.ResponseWriter, request *http.Request) {
	cookieValue := "true"
	host := request.Header.Get(middleware.HTTP1Host)
	domain := host
	if strings.ContainsAny(host, ":") {
		var err error
		domain, _, err = net.SplitHostPort(host)
		if err != nil {
			logging.WithError(err).WithField("host", host).Warning("failed to extract cookie domain from request host")
		}
	}
	a.cookieHandler.SetCookie(writer, a.limitConfig.ExhaustedCookieKey, domain, cookieValue)
}

func (a *AccessInterceptor) DeleteExhaustedCookie(writer http.ResponseWriter) {
	a.cookieHandler.DeleteCookie(writer, a.limitConfig.ExhaustedCookieKey)
}

func (a *AccessInterceptor) Handle(next http.Handler) http.Handler {
	if !a.svc.Enabled() {
		return next
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		tracingCtx, checkSpan := tracing.NewNamedSpan(ctx, "checkAccess")
		wrappedWriter := &statusRecorder{ResponseWriter: writer, status: 0}
		limited := a.Limit(tracingCtx)
		checkSpan.End()
		if limited {
			a.SetExhaustedCookie(wrappedWriter, request)
			http.Error(wrappedWriter, "quota for authenticated requests is exhausted", http.StatusTooManyRequests)
		}
		if !limited && !a.storeOnly {
			a.DeleteExhaustedCookie(wrappedWriter)
		}
		if !limited {
			next.ServeHTTP(wrappedWriter, request)
		}
		tracingCtx, writeSpan := tracing.NewNamedSpan(tracingCtx, "writeAccess")
		defer writeSpan.End()
		requestURL := request.RequestURI
		unescapedURL, err := url.QueryUnescape(requestURL)
		if err != nil {
			logging.WithError(err).WithField("url", requestURL).Warning("failed to unescape request url")
		}
		instance := authz.GetInstance(tracingCtx)
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
