package handler

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/caos/logging"
	"github.com/caos/oidc/pkg/op"
	"github.com/gorilla/csrf"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/api/authz"
	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	auth_repository "github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/auth/repository/eventsourcing"
	"github.com/caos/zitadel/internal/cache"
	cache_config "github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/form"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	_ "github.com/caos/zitadel/internal/ui/login/statik"
	usr_model "github.com/caos/zitadel/internal/user/model"
)

type Login struct {
	endpoint              string
	router                http.Handler
	renderer              *Renderer
	parser                *form.Parser
	command               *command.Commands
	query                 *query.Queries
	staticStorage         static.Storage
	staticCache           cache.Cache
	authRepo              auth_repository.Repository
	baseURL               string
	zitadelURL            string
	oidcAuthCallbackURL   string
	IDPConfigAesCrypto    crypto.EncryptionAlgorithm
	iamDomain             string
	authConnectorVerifier op.JWTProfileVerifier
}

type Config struct {
	BaseURL               string
	OidcAuthCallbackURL   string
	ZitadelURL            string
	Issuer                string
	LanguageCookieName    string
	DefaultLanguage       language.Tag
	CSRF                  CSRF
	UserAgentCookieConfig *middleware.UserAgentCookieConfig
	Cache                 middleware.CacheConfig
	StaticCache           cache_config.CacheConfig
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

func CreateLogin(config Config, command *command.Commands, query *query.Queries, authRepo *eventsourcing.EsRepository, staticStorage static.Storage, systemDefaults systemdefaults.SystemDefaults, localDevMode bool) (*Login, string) {
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
		staticStorage:       staticStorage,
		authRepo:            authRepo,
		IDPConfigAesCrypto:  aesCrypto,
		iamDomain:           systemDefaults.Domain,
	}
	prefix := ""
	if localDevMode {
		prefix = handlerPrefix
	}
	login.staticCache, err = config.StaticCache.Config.NewCache()
	logging.Log("CONFI-dgg31").OnError(err).Panic("unable to create storage cache")

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
	login.renderer = CreateRenderer(prefix, statikFS, staticStorage, config.LanguageCookieName, config.DefaultLanguage)
	login.parser = form.NewParser()
	login.authConnectorVerifier = op.NewJWTProfileVerifier(&jwtProfileStorage{authRepo}, config.Issuer, 30*time.Second, time.Second)
	return login, handlerPrefix
}

type jwtProfileStorage struct {
	repo *eventsourcing.EsRepository
}

func (j jwtProfileStorage) GetKeyByIDAndUserID(ctx context.Context, keyID, userID string) (_ *jose.JSONWebKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
	key, err := j.repo.MachineKeyByID(ctx, keyID)
	if err != nil {
		return nil, err
	}
	if key.AuthIdentifier != userID {
		return nil, errors.ThrowPermissionDenied(nil, "LOGIN-24jm3", "key from different user")
	}
	publicKey, err := crypto.BytesToPublicKey(key.PublicKey)
	if err != nil {
		return nil, err
	}
	return &jose.JSONWebKey{
		KeyID: key.ID,
		Use:   "sig",
		Key:   publicKey,
	}, nil
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
