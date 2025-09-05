package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	zitadel_http "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type instanceInterceptor struct {
	verifier        authz.InstanceVerifier
	externalDomain  string
	ignoredPrefixes []string
	translator      *i18n.Translator
}

func InstanceInterceptor(verifier authz.InstanceVerifier, externalDomain string, translator *i18n.Translator, ignoredPrefixes ...string) *instanceInterceptor {
	return &instanceInterceptor{
		verifier:        verifier,
		externalDomain:  externalDomain,
		ignoredPrefixes: ignoredPrefixes,
		translator:      translator,
	}
}

func (a *instanceInterceptor) Handler(next http.Handler) http.Handler {
	return a.HandlerFunc(next)
}

func (a *instanceInterceptor) HandlerFunc(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := a.setInstanceIfNeeded(r.Context(), r)
		if err == nil {
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		origin := zitadel_http.DomainContext(r.Context())
		logging.WithFields("origin", origin.Origin(), "externalDomain", a.externalDomain).WithError(err).Error("unable to set instance")

		zErr := new(zerrors.ZitadelError)
		if errors.As(err, &zErr) {
			zErr.SetMessage(a.translator.LocalizeFromRequest(r, zErr.GetMessage(), nil))
			http.Error(w, fmt.Sprintf("unable to set instance using origin %s (ExternalDomain is %s): %s", origin, a.externalDomain, zErr), http.StatusNotFound)
			return
		}

		http.Error(w, fmt.Sprintf("unable to set instance using origin %s (ExternalDomain is %s)", origin, a.externalDomain), http.StatusNotFound)
	}
}

func (a *instanceInterceptor) HandlerFuncWithError(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx, err := a.setInstanceIfNeeded(r.Context(), r)
		if err != nil {
			origin := zitadel_http.DomainContext(r.Context())
			logging.WithFields("origin", origin.Origin(), "externalDomain", a.externalDomain).WithError(err).Error("unable to set instance")
			return err
		}

		r = r.WithContext(ctx)
		return next(w, r)
	}
}

func (a *instanceInterceptor) setInstanceIfNeeded(ctx context.Context, r *http.Request) (context.Context, error) {
	for _, prefix := range a.ignoredPrefixes {
		if strings.HasPrefix(r.URL.Path, prefix) {
			return ctx, nil
		}
	}

	return setInstance(ctx, a.verifier)
}

func setInstance(ctx context.Context, verifier authz.InstanceVerifier) (_ context.Context, err error) {
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
