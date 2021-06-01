package handler

import (
	"context"
	"net"
	"net/http"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	auth_repository "github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/form"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
	_ "github.com/caos/zitadel/internal/ui/login/statik"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/gorilla/csrf"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"
)

type Login struct {
	endpoint            string
	router              http.Handler
	renderer            *Renderer
	parser              *form.Parser
	command             *command.Commands
	query               *query.Queries
	authRepo            auth_repository.Repository
	baseURL             string
	zitadelURL          string
	oidcAuthCallbackURL string
	IDPConfigAesCrypto  crypto.EncryptionAlgorithm
	iamDomain           string
}

type Config struct {
	BaseURL               string
	OidcAuthCallbackURL   string
	ZitadelURL            string
	LanguageCookieName    string
	DefaultLanguage       language.Tag
	CSRF                  CSRF
	UserAgentCookieConfig *middleware.UserAgentCookieConfig
	Cache                 middleware.CacheConfig
}

type CSRF struct {
	CookieName  string
	Key         *crypto.KeyConfig
	Development bool
}

const (
	login         = "LOGIN"
	handlerPrefix = "/login"
)

func CreateLogin(config Config, command *command.Commands, query *query.Queries, authRepo *eventsourcing.EsRepository, systemDefaults systemdefaults.SystemDefaults, localDevMode bool) (*Login, string) {
	aesCrypto, err := crypto.NewAESCrypto(systemDefaults.IDPConfigVerificationKey)
	if err != nil {
		logging.Log("HANDL-s90ew").WithError(err).Debug("error create new aes crypto")
	}
	login := &Login{
		oidcAuthCallbackURL: config.OidcAuthCallbackURL,
		baseURL:             config.BaseURL,
		zitadelURL:          config.ZitadelURL,
		command:             command,
		query:               query,
		authRepo:            authRepo,
		IDPConfigAesCrypto:  aesCrypto,
		iamDomain:           systemDefaults.Domain,
	}
	prefix := ""
	if localDevMode {
		prefix = handlerPrefix
	}
	statikFS, err := fs.NewWithNamespace("login")
	logging.Log("CONFI-Ga21f").OnError(err).Panic("unable to create filesystem")

	csrf, err := csrfInterceptor(config.CSRF, login.csrfErrorHandler())
	logging.Log("CONFI-dHR2a").OnError(err).Panic("unable to create csrfInterceptor")
	cache, err := middleware.DefaultCacheInterceptor(EndpointResources, config.Cache.MaxAge.Duration, config.Cache.SharedMaxAge.Duration)
	logging.Log("CONFI-BHq2a").OnError(err).Panic("unable to create cacheInterceptor")
	security := middleware.SecurityHeaders(csp(), login.cspErrorHandler)
	userAgentCookie, err := middleware.NewUserAgentHandler(config.UserAgentCookieConfig, id.SonyFlakeGenerator, localDevMode)
	logging.Log("CONFI-Dvwf2").OnError(err).Panic("unable to create userAgentInterceptor")
	login.router = CreateRouter(login, statikFS, csrf, cache, security, userAgentCookie, middleware.TelemetryHandler(EndpointResources))
	login.renderer = CreateRenderer(prefix, statikFS, config.LanguageCookieName, config.DefaultLanguage)
	login.parser = form.NewParser()
	return login, handlerPrefix
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
	path := "/"
	return csrf.Protect([]byte(csrfKey),
		csrf.Secure(!config.Development),
		csrf.CookieName(http_utils.SetCookiePrefix(config.CookieName, "", path, !config.Development)),
		csrf.Path(path),
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

func (l *Login) getClaimedUserIDsOfOrgDomain(ctx context.Context, orgName string) ([]string, error) {
	users, err := l.authRepo.SearchUsers(ctx, &usr_model.UserSearchRequest{
		Queries: []*usr_model.UserSearchQuery{
			{
				Key:    usr_model.UserSearchKeyPreferredLoginName,
				Method: domain.SearchMethodEndsWithIgnoreCase,
				Value:  domain.NewIAMDomainName(orgName, l.iamDomain),
			},
			{
				Key:    usr_model.UserSearchKeyResourceOwner,
				Method: domain.SearchMethodNotEquals,
				Value:  authz.GetCtxData(ctx).OrgID,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, len(users.Result))
	for i, user := range users.Result {
		userIDs[i] = user.ID
	}
	return userIDs, nil
}

func setContext(ctx context.Context, resourceOwner string) context.Context {
	data := authz.CtxData{
		UserID: login,
		OrgID:  resourceOwner,
	}
	return authz.SetCtxData(ctx, data)
}
