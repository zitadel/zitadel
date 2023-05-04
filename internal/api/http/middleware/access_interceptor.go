package middleware

import (
	"net"
	"net/http"
	"net/url"
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
func NewAccessInterceptor(svc *logstore.Service, cookieHandler *http_utils.CookieHandler, cookieConfig *AccessConfig, storeOnly bool) *AccessInterceptor {
	return &AccessInterceptor{
		svc:           svc,
		cookieHandler: cookieHandler,
		limitConfig:   cookieConfig,
		storeOnly:     storeOnly,
	}
}

func (a *AccessInterceptor) Handle(next http.Handler) http.Handler {
	if !a.svc.Enabled() {
		return next
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		tracingCtx, checkSpan := tracing.NewNamedSpan(ctx, "checkAccess")
		wrappedWriter := &statusRecorder{ResponseWriter: writer, status: 0}
		instance := authz.GetInstance(ctx)
		limit := false
		if !a.storeOnly {
			remaining := a.svc.Limit(tracingCtx, instance.InstanceID())
			limit = remaining != nil && *remaining == 0
		}
		checkSpan.End()
		if limit {
			// Limit can only be true when storeOnly is false, so set the cookie and the response code
			SetExhaustedCookie(a.cookieHandler, wrappedWriter, a.limitConfig, request)
			http.Error(wrappedWriter, "quota for authenticated requests is exhausted", http.StatusTooManyRequests)
		} else {
			if !a.storeOnly {
				// If not limited and not storeOnly, ensure the cookie is deleted
				DeleteExhaustedCookie(a.cookieHandler, wrappedWriter, request, a.limitConfig)
			}
			// Always serve if not limited
			next.ServeHTTP(wrappedWriter, request)
		}
		tracingCtx, writeSpan := tracing.NewNamedSpan(tracingCtx, "writeAccess")
		defer writeSpan.End()
		requestURL := request.RequestURI
		unescapedURL, err := url.QueryUnescape(requestURL)
		if err != nil {
			logging.WithError(err).WithField("url", requestURL).Warning("failed to unescape request url")
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

func SetExhaustedCookie(cookieHandler *http_utils.CookieHandler, writer http.ResponseWriter, cookieConfig *AccessConfig, request *http.Request) {
	cookieValue := "true"
	host := request.Header.Get(middleware.HTTP1Host)
	domain, _, err := net.SplitHostPort(host)
	if err != nil {
		logging.WithError(err).WithField("host", host).Warning("failed to extract cookie domain from request host")
	}
	cookieHandler.SetCookie(writer, cookieConfig.ExhaustedCookieKey, domain, cookieValue)
}

func DeleteExhaustedCookie(cookieHandler *http_utils.CookieHandler, writer http.ResponseWriter, request *http.Request, cookieConfig *AccessConfig) {
	cookieHandler.DeleteCookie(writer, request, cookieConfig.ExhaustedCookieKey)
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
