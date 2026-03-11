package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type AuthInterceptor struct {
	verifier         authz.APITokenVerifier
	authConfig       authz.Config
	systemAuthConfig authz.Config
}

func AuthorizationInterceptor(verifier authz.APITokenVerifier, systemAuthConfig authz.Config, authConfig authz.Config) *AuthInterceptor {
	return &AuthInterceptor{
		verifier:         verifier,
		authConfig:       authConfig,
		systemAuthConfig: systemAuthConfig,
	}
}

func (a *AuthInterceptor) Handler(routePrefix string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return a.HandlerFunc(routePrefix)(next)
	}
}

func (a *AuthInterceptor) HandlerFunc(routePrefix string) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, err := authorize(r, a.verifier, a.systemAuthConfig, a.authConfig, routePrefix)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	}
}

func (a *AuthInterceptor) HandlerFuncWithError(routePrefix string) func(HandlerFuncWithError) HandlerFuncWithError {
	return func(next HandlerFuncWithError) HandlerFuncWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx, err := authorize(r, a.verifier, a.systemAuthConfig, a.authConfig, routePrefix)
			if err != nil {
				return err
			}

			r = r.WithContext(ctx)
			return next(w, r)
		}
	}
}

type httpReq struct{}

func authorize(r *http.Request, verifier authz.APITokenVerifier, systemAuthConfig authz.Config, authConfig authz.Config, routePrefix string) (_ context.Context, err error) {
	ctx := r.Context()

	authOpt, needsToken := checkAuthMethod(r, verifier, routePrefix)
	if !needsToken {
		return ctx, nil
	}
	authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	authToken := http_util.GetAuthorization(r)
	if authToken == "" {
		return nil, zerrors.ThrowUnauthenticated(nil, "AUT-1179", "auth header missing")
	}

	ctxSetter, err := authz.CheckUserAuthorization(authCtx, &httpReq{}, authToken, http_util.GetOrgID(r), "", verifier, systemAuthConfig.RolePermissionMappings, authConfig.RolePermissionMappings, authOpt, r.RequestURI)
	if err != nil {
		return nil, err
	}
	span.End()
	return ctxSetter(ctx), nil
}

func checkAuthMethod(r *http.Request, verifier authz.APITokenVerifier, routePrefix string) (authz.Option, bool) {
	authOpt, needsToken := verifier.CheckAuthMethod(r.Method + ":" + r.RequestURI)
	if needsToken {
		return authOpt, true
	}

	// If the exact path doesn't match, try matching the path template (e.g. /users/{id} instead of /users/123).
	// Since the path template is registered with the sub-router, we need to add the route prefix to it.
	route := mux.CurrentRoute(r)
	if route == nil {
		return authOpt, false
	}
	pathTemplate, err := route.GetPathTemplate()
	if err != nil || pathTemplate == "" {
		return authOpt, false
	}
	return verifier.CheckAuthMethod(r.Method + ":" + routePrefix + pathTemplate)
}
