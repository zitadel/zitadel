package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	zitadel_http "github.com/zitadel/zitadel/v2/internal/api/http"
	"github.com/zitadel/zitadel/v2/internal/i18n"
	"github.com/zitadel/zitadel/v2/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

type instanceInterceptor struct {
	verifier        authz.InstanceVerifier
	externalDomain  string
	ignoredPrefixes []string
	translator      *i18n.Translator
}

func InstanceInterceptor(verifier authz.InstanceVerifier, externalDomain string, ignoredPrefixes ...string) *instanceInterceptor {
	return &instanceInterceptor{
		verifier:        verifier,
		externalDomain:  externalDomain,
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
		origin := zitadel_http.DomainContext(r.Context())
		logging.WithFields("origin", origin.Origin(), "externalDomain", a.externalDomain).WithError(err).Error("unable to set instance")
		zErr := new(zerrors.ZitadelError)
		if errors.As(err, &zErr) {
			zErr.SetMessage(a.translator.LocalizeFromRequest(r, zErr.GetMessage(), nil))
			http.Error(w, fmt.Sprintf("unable to set instance using origin %s (ExternalDomain is %s): %s", origin, a.externalDomain, zErr), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("unable to set instance using origin %s (ExternalDomain is %s)", origin, a.externalDomain), http.StatusNotFound)
		return
	}
	r = r.WithContext(ctx)
	next.ServeHTTP(w, r)
}

func setInstance(r *http.Request, verifier authz.InstanceVerifier) (_ context.Context, err error) {
	ctx := r.Context()
	authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	requestContext := zitadel_http.DomainContext(ctx)
	if requestContext.InstanceHost == "" {
		return nil, zerrors.ThrowNotFound(err, "INST-zWq7X", "Errors.IAM.NotFound")
	}
	instance, err := verifier.InstanceByHost(authCtx, requestContext.InstanceHost, requestContext.PublicHost)
	if err != nil {
		return nil, err
	}
	span.End()
	return authz.WithInstance(ctx, instance), nil
}

func newZitadelTranslator() *i18n.Translator {
	translator, err := i18n.NewZitadelTranslator(language.English)
	logging.OnError(err).Panic("unable to get translator")
	return translator
}
