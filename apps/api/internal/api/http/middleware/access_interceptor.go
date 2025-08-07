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
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type AccessInterceptor struct {
	logstoreSvc   *logstore.Service[*record.AccessLog]
	cookieHandler *http_utils.CookieHandler
	limitConfig   *AccessConfig
	storeOnly     bool
	redirect      string
}

type AccessConfig struct {
	ExhaustedCookieKey    string
	ExhaustedCookieMaxAge time.Duration
}

// NewAccessInterceptor intercepts all requests and stores them to the logstore.
// If storeOnly is false, it also checks if requests are exhausted.
// If requests are exhausted, it also returns http.StatusTooManyRequests or a redirect to the given path and sets a cookie
func NewAccessInterceptor(svc *logstore.Service[*record.AccessLog], cookieHandler *http_utils.CookieHandler, cookieConfig *AccessConfig) *AccessInterceptor {
	return &AccessInterceptor{
		logstoreSvc:   svc,
		cookieHandler: cookieHandler,
		limitConfig:   cookieConfig,
	}
}

func (a *AccessInterceptor) WithoutLimiting() *AccessInterceptor {
	return &AccessInterceptor{
		logstoreSvc:   a.logstoreSvc,
		cookieHandler: a.cookieHandler,
		limitConfig:   a.limitConfig,
		storeOnly:     true,
		redirect:      a.redirect,
	}
}

func (a *AccessInterceptor) WithRedirect(redirect string) *AccessInterceptor {
	return &AccessInterceptor{
		logstoreSvc:   a.logstoreSvc,
		cookieHandler: a.cookieHandler,
		limitConfig:   a.limitConfig,
		storeOnly:     a.storeOnly,
		redirect:      redirect,
	}
}

func (a *AccessInterceptor) AccessService() *logstore.Service[*record.AccessLog] {
	return a.logstoreSvc
}

func (a *AccessInterceptor) Limit(w http.ResponseWriter, r *http.Request, publicAuthPathPrefixes ...string) bool {
	if a.storeOnly {
		return false
	}
	ctx := r.Context()
	instance := authz.GetInstance(ctx)
	var deleteCookie bool
	defer func() {
		if deleteCookie {
			a.DeleteExhaustedCookie(w)
		}
	}()
	if block := instance.Block(); block != nil {
		if *block {
			a.SetExhaustedCookie(w, r)
			return true
		}
		deleteCookie = true
	}
	for _, ignoredPathPrefix := range publicAuthPathPrefixes {
		if strings.HasPrefix(r.RequestURI, ignoredPathPrefix) {
			return false
		}
	}
	remaining := a.logstoreSvc.Limit(ctx, instance.InstanceID())
	if remaining != nil {
		if remaining != nil && *remaining > 0 {
			deleteCookie = true
			return false
		}
		a.SetExhaustedCookie(w, r)
		return true
	}
	return false
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

func (a *AccessInterceptor) HandleWithPublicAuthPathPrefixes(publicPathPrefixes []string) func(next http.Handler) http.Handler {
	return a.handle(publicPathPrefixes...)
}

func (a *AccessInterceptor) Handle(next http.Handler) http.Handler {
	return a.handle()(next)
}

func (a *AccessInterceptor) handle(publicAuthPathPrefixes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			tracingCtx, checkSpan := tracing.NewNamedSpan(ctx, "checkAccessQuota")
			wrappedWriter := &statusRecorder{ResponseWriter: writer, status: 0}
			limited := a.Limit(wrappedWriter, request.WithContext(tracingCtx), publicAuthPathPrefixes...)
			checkSpan.End()
			if limited {
				if a.redirect != "" {
					// The console guides the user when the cookie is set
					http.Redirect(wrappedWriter, request, a.redirect, http.StatusFound)
				} else {
					http.Error(wrappedWriter, "Your ZITADEL instance is blocked.", http.StatusTooManyRequests)
				}
			} else {
				next.ServeHTTP(wrappedWriter, request)
			}
			a.writeLog(tracingCtx, wrappedWriter, writer, request, a.storeOnly)
		})
	}
}

func (a *AccessInterceptor) writeLog(ctx context.Context, wrappedWriter *statusRecorder, writer http.ResponseWriter, request *http.Request, notCountable bool) {
	if !a.logstoreSvc.Enabled() {
		return
	}
	ctx, writeSpan := tracing.NewNamedSpan(ctx, "writeAccess")
	defer writeSpan.End()
	requestURL := request.RequestURI
	unescapedURL, err := url.QueryUnescape(requestURL)
	if err != nil {
		logging.WithError(err).WithField("url", requestURL).Warning("failed to unescape request url")
	}
	instance := authz.GetInstance(ctx)
	domainCtx := http_utils.DomainContext(ctx)
	a.logstoreSvc.Handle(ctx, &record.AccessLog{
		LogDate:         time.Now(),
		Protocol:        record.HTTP,
		RequestURL:      unescapedURL,
		ResponseStatus:  uint32(wrappedWriter.status),
		RequestHeaders:  request.Header,
		ResponseHeaders: writer.Header(),
		InstanceID:      instance.InstanceID(),
		ProjectID:       instance.ProjectID(),
		RequestedDomain: domainCtx.RequestedDomain(),
		RequestedHost:   domainCtx.RequestedHost(),
		NotCountable:    notCountable,
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
