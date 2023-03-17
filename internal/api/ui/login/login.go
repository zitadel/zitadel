package login

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	_ "github.com/zitadel/zitadel/internal/api/ui/login/statik"
	auth_repository "github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/form"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/static"
)

type Login struct {
	endpoint            string
	router              http.Handler
	renderer            *Renderer
	parser              *form.Parser
	command             *command.Commands
	query               *query.Queries
	staticStorage       static.Storage
	authRepo            auth_repository.Repository
	externalSecure      bool
	consolePath         string
	oidcAuthCallbackURL func(context.Context, string) string
	samlAuthCallbackURL func(context.Context, string) string
	idpConfigAlg        crypto.EncryptionAlgorithm
	userCodeAlg         crypto.EncryptionAlgorithm
}

type Config struct {
	LanguageCookieName string
	CSRFCookieName     string
	Cache              middleware.CacheConfig
	AssetCache         middleware.CacheConfig
}

const (
	login                = "LOGIN"
	HandlerPrefix        = "/ui/login"
	DefaultLoggedOutPath = HandlerPrefix + EndpointLogoutDone
)

func CreateLogin(config Config,
	command *command.Commands,
	query *query.Queries,
	authRepo *eventsourcing.EsRepository,
	staticStorage static.Storage,
	consolePath string,
	oidcAuthCallbackURL func(context.Context, string) string,
	samlAuthCallbackURL func(context.Context, string) string,
	externalSecure bool,
	userAgentCookie,
	issuerInterceptor,
	oidcInstanceHandler,
	samlInstanceHandler,
	assetCache,
	accessHandler mux.MiddlewareFunc,
	userCodeAlg crypto.EncryptionAlgorithm,
	idpConfigAlg crypto.EncryptionAlgorithm,
	csrfCookieKey []byte,
) (*Login, error) {
	login := &Login{
		oidcAuthCallbackURL: oidcAuthCallbackURL,
		samlAuthCallbackURL: samlAuthCallbackURL,
		externalSecure:      externalSecure,
		consolePath:         consolePath,
		command:             command,
		query:               query,
		staticStorage:       staticStorage,
		authRepo:            authRepo,
		idpConfigAlg:        idpConfigAlg,
		userCodeAlg:         userCodeAlg,
	}
	statikFS, err := fs.NewWithNamespace("login")
	if err != nil {
		return nil, fmt.Errorf("unable to create filesystem: %w", err)
	}

	csrfInterceptor := createCSRFInterceptor(config.CSRFCookieName, csrfCookieKey, externalSecure, login.csrfErrorHandler())
	cacheInterceptor := createCacheInterceptor(config.Cache.MaxAge, config.Cache.SharedMaxAge, assetCache)
	security := middleware.SecurityHeaders(csp(), login.cspErrorHandler)

	login.router = CreateRouter(login, statikFS, middleware.TelemetryHandler(IgnoreInstanceEndpoints...), oidcInstanceHandler, samlInstanceHandler, csrfInterceptor, cacheInterceptor, security, userAgentCookie, issuerInterceptor, accessHandler)
	login.renderer = CreateRenderer(HandlerPrefix, statikFS, staticStorage, config.LanguageCookieName)
	login.parser = form.NewParser()
	return login, nil
}

func csp() *middleware.CSP {
	csp := middleware.DefaultSCP
	csp.ObjectSrc = middleware.CSPSourceOptsSelf()
	csp.StyleSrc = csp.StyleSrc.AddNonce()
	csp.ScriptSrc = csp.ScriptSrc.AddNonce()
	return &csp
}

func createCSRFInterceptor(cookieName string, csrfCookieKey []byte, externalSecure bool, errorHandler http.Handler) func(http.Handler) http.Handler {
	path := "/"
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, EndpointResources) {
				handler.ServeHTTP(w, r)
				return
			}
			csrf.Protect(csrfCookieKey,
				csrf.Secure(externalSecure),
				csrf.CookieName(http_utils.SetCookiePrefix(cookieName, "", path, externalSecure)),
				csrf.Path(path),
				csrf.ErrorHandler(errorHandler),
			)(handler).ServeHTTP(w, r)
		})
	}
}

func createCacheInterceptor(maxAge, sharedMaxAge time.Duration, assetCache mux.MiddlewareFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, EndpointDynamicResources) {
				assetCache.Middleware(handler).ServeHTTP(w, r)
				return
			}
			if strings.HasPrefix(r.URL.Path, EndpointResources) {
				middleware.AssetsCacheInterceptor(maxAge, sharedMaxAge).Handler(handler).ServeHTTP(w, r)
				return
			}
			middleware.NoCacheInterceptor().Handler(handler).ServeHTTP(w, r)
		})
	}
}

func (l *Login) Handler() http.Handler {
	return l.router
}

func (l *Login) getClaimedUserIDsOfOrgDomain(ctx context.Context, orgName string) ([]string, error) {
	loginName, err := query.NewUserPreferredLoginNameSearchQuery("@"+domain.NewIAMDomainName(orgName, authz.GetInstance(ctx).RequestedDomain()), query.TextEndsWithIgnoreCase)
	if err != nil {
		return nil, err
	}
	users, err := l.query.SearchUsers(ctx, &query.UserSearchQueries{Queries: []query.SearchQuery{loginName}}, false)
	if err != nil {
		return nil, err
	}
	userIDs := make([]string, len(users.Users))
	for i, user := range users.Users {
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

func (l *Login) baseURL(ctx context.Context) string {
	return http_utils.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), l.externalSecure) + HandlerPrefix
}
