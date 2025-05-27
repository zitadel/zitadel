package middleware

import (
	"context"
	"net/http"
	"strings"

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

func (a *AuthInterceptor) Handler(next http.Handler) http.Handler {
	return a.HandlerFunc(next)
}

func (a *AuthInterceptor) HandlerFunc(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := authorize(r, a.verifier, a.systemAuthConfig, a.authConfig)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func (a *AuthInterceptor) HandlerFuncWithError(next HandlerFuncWithError) HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx, err := authorize(r, a.verifier, a.systemAuthConfig, a.authConfig)
		if err != nil {
			return err
		}

		r = r.WithContext(ctx)
		return next(w, r)
	}
}

type httpReq struct{}

func authorize(r *http.Request, verifier authz.APITokenVerifier, systemAuthConfig authz.Config, authConfig authz.Config) (_ context.Context, err error) {
	ctx := r.Context()

	authOpt, needsToken := checkAuthMethod(r, verifier)
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

func checkAuthMethod(r *http.Request, verifier authz.APITokenVerifier) (authz.Option, bool) {
	authOpt, needsToken := verifier.CheckAuthMethod(r.Method + ":" + r.RequestURI)
	if needsToken {
		return authOpt, true
	}

	route := mux.CurrentRoute(r)
	if route == nil {
		return authOpt, false
	}

	pathTemplate, err := route.GetPathTemplate()
	if err != nil || pathTemplate == "" {
		return authOpt, false
	}

	// the path prefix is usually handled in a router in upper layer
	// trim the query and the path of the url to get the correct path prefix
	pathPrefix := r.RequestURI
	if i := strings.Index(pathPrefix, "?"); i != -1 {
		pathPrefix = pathPrefix[0:i]
	}
	pathPrefix = strings.TrimSuffix(pathPrefix, r.URL.Path)

	return verifier.CheckAuthMethod(r.Method + ":" + pathPrefix + pathTemplate)
}
