package middleware

import (
	"bytes"
	"math"
	"net/http"
	"net/url"
	"text/template"
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
	storeOnly     bool
}

type AccessConfig struct {
	ExhaustedCookieKey    string
	ExhaustedCookieValue  string
	ExhaustedCookieMaxAge time.Duration
}

// NewAccessInterceptor intercepts all requests and stores them to the logstore.
// If storeOnly is false, it also checks if requests are exhausted.
// If requests are exhausted, it also returns http.StatusTooManyRequests and sets a cookie
func NewAccessInterceptor(svc *logstore.Service, cookieConfig *AccessConfig, storeOnly bool) *AccessInterceptor {
	return &AccessInterceptor{
		svc: svc,
		cookieHandler: http_utils.NewCookieHandler(
			http_utils.WithUnsecure(),
			http_utils.WithNonHttpOnly(),
			http_utils.WithMaxAge(int(math.Floor(cookieConfig.ExhaustedCookieMaxAge.Seconds()))),
		),
		limitConfig: cookieConfig,
		storeOnly:   storeOnly,
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
		limit := false
		if !a.storeOnly {
			remaining := a.svc.Limit(tracingCtx, instance.InstanceID())
			limit = remaining != nil && *remaining == 0
		}
		if limit {
			// Limit can only be true when storeOnly is false, so set the cookie and the response code
			cookieValue, err := templateCookieValue(a.limitConfig.ExhaustedCookieValue, instance)
			if err != nil {
				// If templating didn't succeed, emit a warning log and just use the plain config
				logging.WithError(err).WithField("value", a.limitConfig.ExhaustedCookieValue).Warning("failed to go template cookie value config")
				err = nil
			}
			a.cookieHandler.SetCookie(wrappedWriter, a.limitConfig.ExhaustedCookieKey, request.Host, cookieValue)
			http.Error(wrappedWriter, "quota for authenticated requests is exhausted", http.StatusTooManyRequests)
		} else {
			if !a.storeOnly {
				// If not limited and not storeOnly, ensure the cookie is deleted
				a.cookieHandler.DeleteCookie(wrappedWriter, request, a.limitConfig.ExhaustedCookieKey)
			}
			// Always serve if not limited
			next.ServeHTTP(wrappedWriter, request)
		}
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

func templateCookieValue(templateableCookieValue string, instance authz.Instance) (string, error) {
	cookieValueTemplate, err := template.New("cookievalue").Parse(templateableCookieValue)
	if err != nil {
		return templateableCookieValue, err
	}
	cookieValue := new(bytes.Buffer)
	if err = cookieValueTemplate.Execute(cookieValue, instance); err != nil {
		return templateableCookieValue, err
	}
	return cookieValue.String(), nil
}
