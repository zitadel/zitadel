package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/rakyll/statik/fs"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	zitadel_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type instanceInterceptor struct {
	verifier        authz.InstanceVerifier
	ignoredPrefixes []string
	translator      *i18n.Translator
}

func InstanceInterceptor(verifier authz.InstanceVerifier, ignoredPrefixes ...string) *instanceInterceptor {
	return &instanceInterceptor{
		verifier:        verifier,
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
	ctx, err := setInstance(r, a.verifier)
	if err != nil {
		caosErr := new(zitadel_errors.NotFoundError)
		if errors.As(err, &caosErr) {
			caosErr.Message = a.translator.LocalizeFromRequest(r, caosErr.GetMessage(), nil)
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	r = r.WithContext(ctx)
	next.ServeHTTP(w, r)
}

func setInstance(r *http.Request, verifier authz.InstanceVerifier) (_ context.Context, err error) {
	ctx := r.Context()
	authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()
	instance, err := verifier.InstanceByDomain(authCtx, http_utils.RequestOriginFromCtx(ctx).Domain)
	if err != nil {
		return nil, err
	}
	span.End()
	return authz.WithInstance(ctx, instance), nil
}

func newZitadelTranslator() *i18n.Translator {
	dir, err := fs.NewWithNamespace("zitadel")
	logging.WithFields("namespace", "zitadel").OnError(err).Panic("unable to get namespace")

	translator, err := i18n.NewTranslator(dir, language.English, "")
	logging.OnError(err).Panic("unable to get translator")
	return translator
}
