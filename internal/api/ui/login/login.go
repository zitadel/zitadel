package login

import (
	"context"
	"net/http"
	"net/url"
	"slices"
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
	DefaultPaths *DefaultPaths
}

type DefaultPaths struct {
	BasePath          *url.URL
	PasswordSetPath   *url.URL
	EmailCodePath     *url.URL
	OTPEmailPath      *url.URL
	PasskeySetPath    *url.URL
	DomainClaimedPath *url.URL
}

func (c *DefaultPaths) defaultBaseURL(ctx context.Context) *url.URL {
	loginV2 := authz.GetInstance(ctx).Features().LoginV2
	// In case login v1 is still active, we don't want to return a base URL, as the templates will not be used
	if !loginV2.Required {
		return nil
	}
	origin := http_utils.DomainContext(ctx).OriginURL()
	if loginV2.BaseURI == nil || loginV2.BaseURI.String() == "" {
		// In case the login v2 is enabled without a custom BaseURI,
		// we use the request origin plus the default base path as the base URL for the templates.
		if c.BasePath == nil {
			return origin
		}
		return origin.ResolveReference(c.BasePath)
	}
	// In case a custom BaseURI is set for login v2, we use it as the base URL for the templates.
	if loginV2.BaseURI.IsAbs() {
		return loginV2.BaseURI
	}
	// If the custom BaseURI is a relative URL, we join the request origin with the custom BaseURI to form the base URL for the templates.
	return origin.ResolveReference(loginV2.BaseURI)
}

// mergeURLs will merge two URLs, where the path is joined and query parameters are combined.
// It uses (*url.URL).JoinPath on the base URL and the second URL's Path field to build the resulting path,
// and then merges the query parameters from both URLs using [mergeQueries], preserving existing values and placeholders.
// Fragments are not modified or combined by this function.
// For example, merging "https://example.com/base" with a URL whose path is "/path/to/resource" will result in "https://example.com/base/path/to/resource".
func mergeURLs(base, path *url.URL) string {
	if base == nil && path == nil {
		return ""
	}
	if base == nil {
		return path.String()
	}
	if path == nil {
		return base.String()
	}
	u := base.JoinPath(path.Path)
	u.RawQuery = mergeQueries(u.Query(), path.Query())
	return u.String()
}

// mergeQueries will merge two sets of query parameters, preserving existing values and placeholders.
// It takes two url.Values, where the first one is considered the base and the second one is merged into it.
// The resulting query string will contain all unique key-value pairs from both sets of query parameters.
// If a value contains placeholders in the format "{{placeholder}}", they are preserved in the final query string without being URL-encoded.
func mergeQueries(base, path url.Values) string {
	placeholders := map[string]string{}
	for _, value := range base {
		for _, v := range value {
			if strings.Contains(v, "{{") && strings.Contains(v, "}}") {
				placeholders[url.QueryEscape(v)] = v
			}
		}
	}
	for key, value := range path {
		baseValues, ok := base[key]
		for _, v := range value {
			if !ok || !slices.Contains(baseValues, v) {
				base.Add(key, v)
			}
			if strings.Contains(v, "{{") && strings.Contains(v, "}}") {
				placeholders[url.QueryEscape(v)] = v
			}
		}
	}
	raw := base.Encode()
	for esc, orig := range placeholders {
		raw = strings.ReplaceAll(raw, esc, orig)
	}
	return raw
}

func (c *DefaultPaths) DefaultEmailCodeURLTemplate(ctx context.Context) string {
	basePath := c.defaultBaseURL(ctx)
	if basePath == nil {
		return ""
	}
	return mergeURLs(basePath, c.EmailCodePath)
}

func (c *DefaultPaths) DefaultPasswordSetURLTemplate(ctx context.Context) string {
	basePath := c.defaultBaseURL(ctx)
	if basePath == nil {
		return ""
	}
	return mergeURLs(basePath, c.PasswordSetPath)
}

func (c *DefaultPaths) DefaultPasskeySetURLTemplate(ctx context.Context) string {
	basePath := c.defaultBaseURL(ctx)
	if basePath == nil {
		return ""
	}
	return mergeURLs(basePath, c.PasskeySetPath)
}

func (c *DefaultPaths) DefaultDomainClaimedURLTemplate(ctx context.Context) string {
	basePath := c.defaultBaseURL(ctx)
	if basePath == nil {
		return ""
	}
	return mergeURLs(basePath, c.DomainClaimedPath)
}

func (c *DefaultPaths) DefaultOTPEmailURLTemplate(origin *url.URL) string {
	if c.BasePath == nil {
		return mergeURLs(origin, c.OTPEmailPath)
	}
	return mergeURLs(origin.ResolveReference(c.BasePath), c.OTPEmailPath)
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

	login.router = CreateRouter(login,
		middleware.TraceHandler(IgnoreInstanceEndpoints...),
		middleware.LogHandler("login_v1", IgnoreInstanceEndpoints...),
		oidcInstanceHandler,
		samlInstanceHandler,
		csrfInterceptor,
		cacheInterceptor,
		security,
		userAgentCookie,
		issuerInterceptor,
		accessHandler,
	)
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
