package handler

import (
	"context"
	"net"
	"net/http"

	"github.com/caos/logging"
	"github.com/gorilla/csrf"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/form"
	_ "github.com/caos/zitadel/internal/ui/login/statik"
)

type Login struct {
	endpoint            string
	router              http.Handler
	renderer            *Renderer
	parser              *form.Parser
	authRepo            *eventsourcing.EsRepository
	zitadelURL          string
	oidcAuthCallbackURL string
}

type Config struct {
	OidcAuthCallbackURL string
	ZitadelURL          string
	LanguageCookieName  string
	DefaultLanguage     language.Tag
	CSRF                CSRF
	Cache               middleware.CacheConfig
}

type CSRF struct {
	CookieName  string
	Key         *crypto.KeyConfig
	Development bool
}

const (
	login = "LOGIN"
)

func CreateLogin(config Config, authRepo *eventsourcing.EsRepository, prefix string) *Login {
	login := &Login{
		oidcAuthCallbackURL: config.OidcAuthCallbackURL,
		zitadelURL:          config.ZitadelURL,
		authRepo:            authRepo,
	}
	statikFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-Ga21f").OnError(err).Panic("unable to create filesystem")

	csrf, err := csrfInterceptor(config.CSRF, login.csrfErrorHandler())
	logging.Log("CONFI-dHR2a").OnError(err).Panic("unable to create csrfInterceptor")
	cache, err := middleware.DefaultCacheInterceptor(EndpointResources, config.Cache.MaxAge.Duration, config.Cache.SharedMaxAge.Duration)
	logging.Log("CONFI-BHq2a").OnError(err).Panic("unable to create cacheInterceptor")
	security := middleware.SecurityHeaders(csp(), login.cspErrorHandler)
	login.router = CreateRouter(login, statikFS, csrf, cache, security)
	login.renderer = CreateRenderer(prefix, statikFS, config.LanguageCookieName, config.DefaultLanguage)
	login.parser = form.NewParser()
	return login
}

func csp() *middleware.CSP {
	csp := middleware.DefaultSCP
	csp.ObjectSrc = middleware.CSPSourceOptsSelf()
	csp.StyleSrc = csp.StyleSrc.AddNonce()
	csp.ScriptSrc = csp.ScriptSrc.AddNonce()
	return &csp
}

func csrfInterceptor(config CSRF, errorHandler http.Handler) (func(http.Handler) http.Handler, error) {
	csrfKey, err := crypto.LoadKey(config.Key, config.Key.EncryptionKeyID)
	if err != nil {
		return nil, err
	}
	return csrf.Protect([]byte(csrfKey),
		csrf.Secure(!config.Development),
		csrf.CookieName(config.CookieName),
		csrf.ErrorHandler(errorHandler),
	), nil
}

func (l *Login) Handler() http.Handler {
	return l.router
}

func (l *Login) Listen(ctx context.Context) {
	if l.endpoint == "" {
		l.endpoint = ":80"
	} else {
		l.endpoint = ":" + l.endpoint
	}

	defer logging.LogWithFields("APP-xUZof", "port", l.endpoint).Info("html is listening")
	httpListener, err := net.Listen("tcp", l.endpoint)
	logging.Log("CONFI-W5q2O").OnError(err).Panic("unable to start listener")

	httpServer := &http.Server{
		Handler: l.router,
	}

	go func() {
		<-ctx.Done()
		if err = httpServer.Shutdown(ctx); err != nil {
			logging.Log("APP-mJKTv").WithError(err)
		}
	}()

	go func() {
		err := httpServer.Serve(httpListener)
		logging.Log("APP-oSklt").OnError(err).Panic("unable to start listener")
	}()
}

func setContext(ctx context.Context, resourceOwner string) context.Context {
	data := authz.CtxData{
		UserID: login,
		OrgID:  resourceOwner,
	}
	return authz.SetCtxData(ctx, data)
}
