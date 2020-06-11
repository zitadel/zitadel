package handler

import (
	"context"
	"net"
	"net/http"

	"github.com/caos/logging"
	"github.com/gorilla/csrf"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/form"

	_ "github.com/caos/zitadel/internal/login/statik"
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
	Port                string
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

func StartLogin(ctx context.Context, config Config, authRepo *eventsourcing.EsRepository) {
	login := &Login{
		endpoint:            config.Port,
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
	security := middleware.SecurityHeaders(csp(config.OidcAuthCallbackURL))
	login.router = CreateRouter(login, statikFS, csrf, cache, security)
	login.renderer = CreateRenderer(statikFS, config.LanguageCookieName, config.DefaultLanguage)
	login.parser = form.NewParser()
	login.Listen(ctx)
}

func csp(callback string) *middleware.CSP {
	csp := middleware.DefaultSCP
	csp.StyleSrc.AddNonce()
	csp.ScriptSrc.AddNonce()
	csp.FormAction.AddHost(callback)
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
	data := auth.CtxData{
		UserID: login,
		OrgID:  resourceOwner,
	}
	return auth.SetCtxData(ctx, data)
}
