package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type instanceInterceptor struct {
	verifier   authz.InstanceVerifier
	headerName string
}

func InstanceInterceptor(verifier authz.InstanceVerifier, headerName string) *instanceInterceptor {
	return &instanceInterceptor{
		verifier:   verifier,
		headerName: headerName,
	}
}

func (a *instanceInterceptor) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, err := setInstance(r, a.verifier, a.headerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (a *instanceInterceptor) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := setInstance(r, a.verifier, a.headerName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func setInstance(r *http.Request, verifier authz.InstanceVerifier, headerName string) (_ context.Context, err error) {
	ctx := r.Context()

	authCtx, span := tracing.NewServerInterceptorSpan(ctx)
	defer func() { span.EndWithError(err) }()

	host, err := getHost(r, headerName)
	if err != nil {
		return nil, err
	}

	instance, err := verifier.InstanceByHost(authCtx, host)
	if err != nil {
		return nil, err
	}
	span.End()
	return authz.WithInstance(ctx, instance), nil
}

func getHost(r *http.Request, headerName string) (string, error) {
	host := r.Host
	if headerName != "host" {
		host = r.Header.Get(headerName)
	}
	if host == "" {
		return "", fmt.Errorf("host header `%s` not found", headerName)
	}
	return host, nil
}
