package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/caos/zitadel/internal/api/authz"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type AuthInterceptor struct {
	verifier   *authz.TokenVerifier
	authConfig authz.Config
}

func AuthorizationInterceptor(verifier *authz.TokenVerifier, authConfig authz.Config) *AuthInterceptor {
	return &AuthInterceptor{
		verifier:   verifier,
		authConfig: authConfig,
	}
}

func (a *AuthInterceptor) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := authorize(r, a.verifier, a.authConfig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

type httpReq struct{}

func authorize(r *http.Request, verifier *authz.TokenVerifier, authConfig authz.Config) (_ context.Context, err error) {
	ctx := r.Context()
	authOpt, needsToken := verifier.CheckAuthMethod(r.Method + ":" + r.RequestURI)
	if !needsToken {
		return ctx, nil
	}
	authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	authToken := http_util.GetAuthorization(r)
	if authToken == "" {
		return nil, errors.New("auth header missing")
	}

	ctxSetter, err := authz.CheckUserAuthorization(authCtx, &httpReq{}, authToken, http_util.GetOrgID(r), verifier, authConfig, authOpt, r.RequestURI) //TODO: permission
	if err != nil {
		return nil, err
	}
	span.End()
	return ctxSetter(ctx), nil
}
