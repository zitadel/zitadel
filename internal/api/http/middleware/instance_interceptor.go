package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rakyll/statik/fs"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	caos_errors "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type instanceInterceptor struct {
	verifier         authz.InstanceVerifier
	virtualInstances bool
	headerName       string
	ignoredPrefixes  []string
	translator       *i18n.Translator
}

func InstanceInterceptor(verifier authz.InstanceVerifier, virtualInstances bool, headerName string, ignoredPrefixes ...string) *instanceInterceptor {
	return &instanceInterceptor{
		verifier:         verifier,
		virtualInstances: virtualInstances,
		headerName:       headerName,
		ignoredPrefixes:  ignoredPrefixes,
		translator:       newZitadelTranslator(),
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
	var (
		instance authz.Instance
		err      error
	)
	if a.virtualInstances {
		instance, err = queryInstanceByHost(r, a.verifier, a.headerName)
	} else {
		instance, err = queryFirstInstance(r, a.verifier)
	}
	if err != nil {
		caosErr := new(caos_errors.NotFoundError)
		if errors.As(err, &caosErr) {
			caosErr.Message = a.translator.LocalizeFromRequest(r, caosErr.GetMessage(), nil)
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	r = r.WithContext(authz.WithInstance(r.Context(), instance))
	next.ServeHTTP(w, r)
}

func queryFirstInstance(r *http.Request, verifier authz.InstanceVerifier) (instance authz.Instance, err error) {
	ctx, span := tracing.NewServerInterceptorSpan(r.Context())
	defer func() { span.EndWithError(err) }()
	return verifier.FirstInstance(ctx)
}

func queryInstanceByHost(r *http.Request, verifier authz.InstanceVerifier, headerName string) (instance authz.Instance, err error) {
	ctx, span := tracing.NewServerInterceptorSpan(r.Context())
	defer func() { span.EndWithError(err) }()
	host, err := HostFromRequest(r, headerName)
	if err != nil {
		return nil, err
	}
	return verifier.InstanceByHost(ctx, host)
}

func HostFromRequest(r *http.Request, headerName string) (string, error) {
	host := r.Host
	if headerName != "host" {
		host = r.Header.Get(headerName)
	}
	if host == "" {
		return "", fmt.Errorf("host header `%s` not found", headerName)
	}
	return host, nil
}

func newZitadelTranslator() *i18n.Translator {
	dir, err := fs.NewWithNamespace("zitadel")
	logging.WithFields("namespace", "zitadel").OnError(err).Panic("unable to get namespace")

	translator, err := i18n.NewTranslator(dir, language.English, "")
	logging.OnError(err).Panic("unable to get translator")
	return translator
}
