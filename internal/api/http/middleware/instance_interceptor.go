package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type instanceInterceptor struct {
	verifier        authz.InstanceVerifier
	headerName      string
	ignoredPrefixes []string
	translator      *i18n.Translator
}

func InstanceInterceptor(verifier authz.InstanceVerifier, headerName string, ignoredPrefixes ...string) *instanceInterceptor {
	return &instanceInterceptor{
		verifier:        verifier,
		headerName:      headerName,
		ignoredPrefixes: ignoredPrefixes,
		translator:      newZitadelTranslator(),
	}
}

func (a *instanceInterceptor) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.handleInstance(w, r, next)
	})
}

func (a *instanceInterceptor) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.handleInstance(w, r, next)
	}
}

func (a *instanceInterceptor) handleInstance(w http.ResponseWriter, r *http.Request, next http.Handler) {
	for _, prefix := range a.ignoredPrefixes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			next.ServeHTTP(w, r)
			return
		}
	}
	ctx, err := setInstance(r, a.verifier, a.headerName)
	if err != nil {
		caosErr := new(zerrors.NotFoundError)
		if errors.As(err, &caosErr) {
			caosErr.Message = a.translator.LocalizeFromRequest(r, caosErr.GetMessage(), nil)
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	r = r.WithContext(ctx)
	next.ServeHTTP(w, r)
}

func setInstance(r *http.Request, verifier authz.InstanceVerifier, headerName string) (_ context.Context, err error) {
	ctx := r.Context()

	authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	host, err := HostFromRequest(r, headerName)
	if err != nil {
		return nil, zerrors.ThrowNotFound(err, "INST-zWq7X", "Errors.Instance.NotFound")
	}

	instance, err := verifier.InstanceByHost(authCtx, host)
	if err != nil {
		return nil, err
	}
	span.End()
	return authz.WithInstance(ctx, instance), nil
}

func HostFromRequest(r *http.Request, headerName string) (host string, err error) {
	if headerName != "host" {
		return hostFromSpecialHeader(r, headerName)
	}
	return hostFromOrigin(r.Context())
}

func hostFromSpecialHeader(r *http.Request, headerName string) (host string, err error) {
	host = r.Header.Get(headerName)
	if host == "" {
		return "", fmt.Errorf("host header `%s` not found", headerName)
	}
	return host, nil
}

func hostFromOrigin(ctx context.Context) (host string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("invalid origin: %w", err)
		}
	}()
	origin := zitadel_http.ComposedOrigin(ctx)
	u, err := url.Parse(origin)
	if err != nil {
		return "", err
	}
	host = u.Host
	if host == "" {
		err = errors.New("empty host")
	}
	return host, err
}

func newZitadelTranslator() *i18n.Translator {
	translator, err := i18n.NewZitadelTranslator(language.English)
	logging.OnError(err).Panic("unable to get translator")
	return translator
}
