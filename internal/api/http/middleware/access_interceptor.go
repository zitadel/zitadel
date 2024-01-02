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
	"github.com/zitadel/zitadel/internal/api/limits"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type AccessInterceptor struct {
	logstoreSvc   *logstore.Service[*record.AccessLog]
	cookieHandler *http_utils.CookieHandler
	limitConfig   *AccessConfig
	storeOnly     bool
	limitsLoader  *limits.Loader
}

type AccessConfig struct {
	ExhaustedCookieKey    string
	ExhaustedCookieMaxAge time.Duration
}

// NewAccessInterceptor intercepts all requests and stores them to the logstore.
// If storeOnly is false, it also checks if requests are exhausted.
// If requests are exhausted, it also returns http.StatusTooManyRequests and sets a cookie
func NewAccessInterceptor(svc *logstore.Service[*record.AccessLog], limitsLoader *limits.Loader, cookieHandler *http_utils.CookieHandler, cookieConfig *AccessConfig) *AccessInterceptor {
	return &AccessInterceptor{
		logstoreSvc:   svc,
		cookieHandler: cookieHandler,
		limitConfig:   cookieConfig,
		limitsLoader:  limitsLoader,
	}
}

func (a *AccessInterceptor) WithoutLimiting() *AccessInterceptor {
	return &AccessInterceptor{
		logstoreSvc:   a.logstoreSvc,
		cookieHandler: a.cookieHandler,
		limitConfig:   a.limitConfig,
		storeOnly:     true,
	}
}

func (a *AccessInterceptor) AccessService() *logstore.Service[*record.AccessLog] {
	return a.logstoreSvc
}

func (a *AccessInterceptor) Limit(ctx context.Context) (context.Context, bool) {
	if a.storeOnly {
		return ctx, false
	}
	instanceID := authz.GetInstance(ctx).InstanceID()
	ctx, l := a.limitsLoader.Load(ctx, instanceID)
	if l.Block != nil && *l.Block {
		return ctx, true
	}
	if !a.logstoreSvc.Enabled() {
		return ctx, false
	}
	remaining := a.logstoreSvc.Limit(ctx, instanceID)
	return ctx, remaining != nil && *remaining <= 0
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

func (a *AccessInterceptor) HandleIgnorePathPrefixes(ignoredPathPrefixes []string) func(next http.Handler) http.Handler {
	return a.handle(ignoredPathPrefixes...)
}

func (a *AccessInterceptor) Handle(next http.Handler) http.Handler {
	return a.handle()(next)
}

func (a *AccessInterceptor) handle(ignoredPathPrefixes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			ctx := request.Context()
			tracingCtx, checkSpan := tracing.NewNamedSpan(ctx, "checkAccessQuota")
			wrappedWriter := &statusRecorder{ResponseWriter: writer, status: 0}
			for _, ignoredPathPrefix := range ignoredPathPrefixes {
				if !strings.HasPrefix(request.RequestURI, ignoredPathPrefix) {
					continue
				}
				checkSpan.End()
				next.ServeHTTP(wrappedWriter, request)
				a.writeLog(tracingCtx, wrappedWriter, writer, request, true)
				return
			}
			ctx, limited := a.Limit(tracingCtx)
			request = request.WithContext(ctx)
			checkSpan.End()
			if limited {
				a.SetExhaustedCookie(wrappedWriter, request)
				if strings.HasPrefix(request.RequestURI, "/ui/login") {
					// The console guides the user when the cookie is set
					http.Redirect(wrappedWriter, request, "/ui/console", http.StatusPermanentRedirect)
				} else {
					http.Error(wrappedWriter, "quota for authenticated requests is exhausted", http.StatusTooManyRequests)
				}
			}
			if !limited && !a.storeOnly {
				a.DeleteExhaustedCookie(wrappedWriter)
			}
			if !limited {
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
	a.logstoreSvc.Handle(ctx, &record.AccessLog{
		LogDate:         time.Now(),
		Protocol:        record.HTTP,
		RequestURL:      unescapedURL,
		ResponseStatus:  uint32(wrappedWriter.status),
		RequestHeaders:  request.Header,
		ResponseHeaders: writer.Header(),
		InstanceID:      instance.InstanceID(),
		ProjectID:       instance.ProjectID(),
		RequestedDomain: instance.RequestedDomain(),
		RequestedHost:   instance.RequestedHost(),
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
