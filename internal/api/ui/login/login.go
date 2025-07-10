package login

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	_ "github.com/zitadel/zitadel/internal/api/ui/login/statik"
	auth_repository "github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/auth/repository/eventsourcing"
	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/connector"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/domain/federatedlogout"
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
	caches              *Caches
}

type Config struct {
	LanguageCookieName string
	CSRFCookieName     string
	Cache              middleware.CacheConfig
	AssetCache         middleware.CacheConfig

	// LoginV2
	DefaultOTPEmailURLV2 string
}

const (
	login                = "LOGIN"
	HandlerPrefix        = "/ui/login"
	DefaultLoggedOutPath = HandlerPrefix + EndpointLogoutDone
)

func CreateLogin(
	config Config,
	command *command.Commands,
	query *query.Queries,
	authRepo *eventsourcing.EsRepository,
	staticStorage static.Storage,
	consolePath string,
	oidcAuthCallbackURL, samlAuthCallbackURL func(context.Context, string) string,
	externalSecure bool,
	userAgentCookie, issuerInterceptor, oidcInstanceHandler, samlInstanceHandler, assetCache, accessHandler mux.MiddlewareFunc,
	userCodeAlg, idpConfigAlg crypto.EncryptionAlgorithm,
	csrfCookieKey []byte,
	cacheConnectors connector.Connectors,
	federateLogoutCache cache.Cache[federatedlogout.Index, string, *federatedlogout.FederatedLogout],
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
	csrfInterceptor := createCSRFInterceptor(config.CSRFCookieName, csrfCookieKey, externalSecure, login.csrfErrorHandler())
	cacheInterceptor := createCacheInterceptor(config.Cache.MaxAge, config.Cache.SharedMaxAge, assetCache)
	security := middleware.SecurityHeaders(csp(), login.cspErrorHandler)

	login.router = CreateRouter(login, middleware.TelemetryHandler(IgnoreInstanceEndpoints...), oidcInstanceHandler, samlInstanceHandler, csrfInterceptor, cacheInterceptor, security, userAgentCookie, issuerInterceptor, accessHandler)
	login.renderer = CreateRenderer(HandlerPrefix, staticStorage, config.LanguageCookieName)
	login.parser = form.NewParser()

	var err error
	login.caches, err = startCaches(context.Background(), cacheConnectors, federateLogoutCache)
	if err != nil {
		return nil, err
	}
	return login, nil
}

func csp() *middleware.CSP {
	csp := middleware.DefaultSCP
	csp.ObjectSrc = middleware.CSPSourceOptsSelf()
	csp.StyleSrc = csp.StyleSrc.AddNonce()
	csp.ScriptSrc = csp.ScriptSrc.AddNonce().
		// SAML POST ACS
		AddHash("sha256", "AjPdJSbZmeWHnEc5ykvJFay8FTWeTeRbs9dutfZ0HqE=").
		// SAML POST SLO
		AddHash("sha256", "4Su6mBWzEIFnH4pAGMOuaeBrstwJN4Z3pq/s1Kn4/KQ=")
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
			// ignore form post callback
			// it will redirect to the "normal" callback, where the cookie is set again
			if (r.URL.Path == EndpointExternalLoginCallbackFormPost || r.URL.Path == EndpointSAMLACS) && r.Method == http.MethodPost {
				handler.ServeHTTP(w, r)
				return
			}
			// by default we use SameSite Lax and the externalSecure (TLS) for the secure flag
			sameSiteMode := csrf.SameSiteLaxMode
			secureOnly := externalSecure
			instance := authz.GetInstance(r.Context())
			// in case of `allow iframe`...
			if len(instance.SecurityPolicyAllowedOrigins()) > 0 {
				// ... we need to change to SameSite none ...
				sameSiteMode = csrf.SameSiteNoneMode
				// ... and since SameSite none requires the secure flag, we'll set it for TLS and for localhost
				// (regardless of the TLS / externalSecure settings)
				secureOnly = externalSecure || http_utils.DomainContext(r.Context()).RequestedDomain() == "localhost"
			}
			csrf.Protect(csrfCookieKey,
				csrf.Secure(secureOnly),
				csrf.CookieName(http_utils.SetCookiePrefix(cookieName, externalSecure, http_utils.PrefixHost)),
				csrf.Path(path),
				csrf.ErrorHandler(errorHandler),
				csrf.SameSite(sameSiteMode),
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
	orgDomain, err := domain.NewIAMDomainName(orgName, http_utils.DomainContext(ctx).RequestedDomain())
	if err != nil {
		return nil, err
	}
	return l.query.SearchClaimedUserIDsOfOrgDomain(ctx, orgDomain, "")
}

func setContext(ctx context.Context, resourceOwner string) context.Context {
	data := authz.CtxData{
		UserID: login,
		OrgID:  resourceOwner,
	}
	return authz.SetCtxData(ctx, data)
}

func setUserContext(ctx context.Context, userID, resourceOwner string) context.Context {
	data := authz.CtxData{
		UserID: userID,
		OrgID:  resourceOwner,
	}
	return authz.SetCtxData(ctx, data)
}

func (l *Login) baseURL(ctx context.Context) string {
	return http_utils.DomainContext(ctx).Origin() + HandlerPrefix
}

type Caches struct {
	idpFormCallbacks cache.Cache[idpFormCallbackIndex, string, *idpFormCallback]
	federatedLogouts cache.Cache[federatedlogout.Index, string, *federatedlogout.FederatedLogout]
}

func startCaches(background context.Context, connectors connector.Connectors, federateLogoutCache cache.Cache[federatedlogout.Index, string, *federatedlogout.FederatedLogout]) (_ *Caches, err error) {
	caches := new(Caches)
	caches.idpFormCallbacks, err = connector.StartCache[idpFormCallbackIndex, string, *idpFormCallback](background, []idpFormCallbackIndex{idpFormCallbackIndexRequestID}, cache.PurposeIdPFormCallback, connectors.Config.IdPFormCallbacks, connectors)
	if err != nil {
		return nil, err
	}
	caches.federatedLogouts = federateLogoutCache
	return caches, nil
}

type idpFormCallbackIndex int

const (
	idpFormCallbackIndexUnspecified idpFormCallbackIndex = iota
	idpFormCallbackIndexRequestID
)

type idpFormCallback struct {
	InstanceID string
	State      string
	Form       url.Values
}

// Keys implements cache.Entry
func (c *idpFormCallback) Keys(i idpFormCallbackIndex) []string {
	if i == idpFormCallbackIndexRequestID {
		return []string{idpFormCallbackKey(c.InstanceID, c.State)}
	}
	return nil
}

func idpFormCallbackKey(instanceID, state string) string {
	return instanceID + "-" + state
}
