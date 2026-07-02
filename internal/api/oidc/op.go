package oidc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/metrics"
	"github.com/zitadel/zitadel/internal/api/assets"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain/federatedlogout"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/handlers"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	CodeMethodS256                    bool
	AuthMethodPost                    bool
	AuthMethodPrivateKeyJWT           bool
	GrantTypeRefreshToken             bool
	RequestObjectSupported            bool
	DefaultAccessTokenLifetime        time.Duration
	DefaultIdTokenLifetime            time.Duration
	DefaultRefreshTokenIdleExpiration time.Duration
	DefaultRefreshTokenExpiration     time.Duration
	JWKSCacheControlMaxAge            time.Duration
	CustomEndpoints                   *EndpointConfig
	DeviceAuth                        *DeviceAuthorizationConfig
	DefaultLoginURLV2                 string
	DefaultLogoutURLV2                string
	PublicKeyCacheMaxAge              time.Duration
	DefaultBackChannelLogoutLifetime  time.Duration
	BackChannelLogout                 handlers.BackChannelLogoutWorkerConfig
	DynamicClientRegistration         DynamicClientRegistrationConfig
}

// DynamicClientRegistrationConfig configures the OAuth 2.0 Dynamic Client Registration
// endpoint (RFC 7591). The endpoint is only served and advertised when the
// oidc_dynamic_client_registration instance feature is enabled.
type DynamicClientRegistrationConfig struct {
	// AllowUnauthenticated enables open registration, where clients may register
	// without an initial access token, as required by the Model Context Protocol (MCP)
	// flow. Registered clients are homed in a dedicated, auto-provisioned project in the
	// instance's default organization. When disabled (the default), registration
	// requires a valid access token (RFC 7591 §3 initial access token) and the client is
	// homed in the token's organization.
	AllowUnauthenticated bool
}

// BackChannelLogoutConfig returns the BackChannelLogoutWorkerConfig and takes the deprecated TokenLifetime into account.
func (c *Config) BackChannelLogoutConfig() *handlers.BackChannelLogoutWorkerConfig {
	if c.DefaultBackChannelLogoutLifetime == 0 {
		return &c.BackChannelLogout
	}
	c.BackChannelLogout.TokenLifetime = c.DefaultBackChannelLogoutLifetime
	return &c.BackChannelLogout
}

type EndpointConfig struct {
	Auth          *Endpoint
	Token         *Endpoint
	Introspection *Endpoint
	Userinfo      *Endpoint
	Revocation    *Endpoint
	EndSession    *Endpoint
	Keys          *Endpoint
	DeviceAuth    *Endpoint
	Registration  *Endpoint
}

type Endpoint struct {
	Path string
	URL  string
}

type OPStorage struct {
	repo                              repository.Repository
	command                           *command.Commands
	query                             *query.Queries
	eventstore                        *eventstore.Eventstore
	defaultLoginURL                   string
	defaultLoginURLV2                 string
	defaultLogoutURLV2                string
	defaultAccessTokenLifetime        time.Duration
	defaultIdTokenLifetime            time.Duration
	defaultRefreshTokenIdleExpiration time.Duration
	defaultRefreshTokenExpiration     time.Duration
	authAlg                           crypto.AuthAlgorithm
	assetAPIPrefix                    func(ctx context.Context) string
	contextToIssuer                   func(context.Context) string
	federateLogoutCache               cache.Cache[federatedlogout.Index, string, *federatedlogout.FederatedLogout]
	clientIDMetadataResolver          *clientIDMetadataResolver
}

// Provider is used to overload certain [op.Provider] methods
type Provider struct {
	*op.Provider
	accessTokenKeySet oidc.KeySet
	idTokenHintKeySet oidc.KeySet
}

// IDTokenHintVerifier configures a Verifier and supported signing algorithms based on the Web Key feature in the context.
func (o *Provider) IDTokenHintVerifier(ctx context.Context) *op.IDTokenHintVerifier {
	return op.NewIDTokenHintVerifier(op.IssuerFromContext(ctx), o.idTokenHintKeySet, op.WithSupportedIDTokenHintSigningAlgorithms(
		supportedSigningAlgs()...,
	))
}

// AccessTokenVerifier configures a Verifier and supported signing algorithms based on the Web Key feature in the context.
func (o *Provider) AccessTokenVerifier(ctx context.Context) *op.AccessTokenVerifier {
	return op.NewAccessTokenVerifier(op.IssuerFromContext(ctx), o.accessTokenKeySet, op.WithSupportedAccessTokenSigningAlgorithms(
		supportedSigningAlgs()...,
	))
}

func NewServer(
	ctx context.Context,
	config Config,
	defaultLogoutRedirectURI string,
	externalSecure bool,
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	authAlg crypto.AuthAlgorithm,
	targetEncryptionAlgorithm crypto.EncryptionAlgorithm,
	cryptoKey []byte,
	es *eventstore.Eventstore,
	userAgentCookie, instanceHandler func(http.Handler) http.Handler,
	accessHandler *middleware.AccessInterceptor,
	fallbackLogger *slog.Logger,
	hashConfig crypto.HashConfig,
	federatedLogoutCache cache.Cache[federatedlogout.Index, string, *federatedlogout.FederatedLogout],
	clientIDMetadataDocumentCache cache.Cache[clientIDMetadataCacheIndex, string, *clientIDMetadataCacheEntry],
	httpClient *http.Client,
) (*Server, error) {
	opConfig, err := createOPConfig(config, defaultLogoutRedirectURI, cryptoKey)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-EGrqd", "cannot create op config: %w")
	}
	clientIDMetadataResolver := newClientIDMetadataResolver(httpClient, clientIDMetadataDocumentCache, config.DefaultAccessTokenLifetime, config.DefaultIdTokenLifetime)
	storage := newStorage(config, command, query, repo, authAlg, es, ContextToIssuer, federatedLogoutCache, clientIDMetadataResolver)
	keyCache := newPublicKeyCache(ctx, config.PublicKeyCacheMaxAge, queryKeyFunc(query))
	accessTokenKeySet := newOidcKeySet(keyCache, withKeyExpiryCheck(true))
	idTokenHintKeySet := newOidcKeySet(keyCache)

	alg := op.NewAES256GCMCrypto(opConfig.CryptoKey, "")
	if authAlg.LegacyTokenEnabled() {
		alg = op.NewCompositeCrypto(
			alg,
			[]op.Decrypter{
				alg,
				op.NewAESCrypto(opConfig.CryptoKey),
			},
		)
	}
	options := []op.Option{
		op.WithCrypto(alg),
	}
	if !externalSecure {
		options = append(options, op.WithAllowInsecure())
	}
	provider, err := op.NewProvider(
		opConfig,
		storage,
		IssuerFromContext,
		options...,
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-DAtg3", "cannot create provider")
	}
	hasher, err := hashConfig.NewHasher()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-Aij4e", "cannot create secret hasher")
	}
	server := &Server{
		LegacyServer: op.NewLegacyServer(&Provider{
			Provider:          provider,
			accessTokenKeySet: accessTokenKeySet,
			idTokenHintKeySet: idTokenHintKeySet,
		}, endpoints(config.CustomEndpoints)),
		repo:                            repo,
		query:                           query,
		command:                         command,
		accessTokenKeySet:               accessTokenKeySet,
		idTokenHintKeySet:               idTokenHintKeySet,
		defaultLoginURL:                 fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID),
		defaultLoginURLV2:               config.DefaultLoginURLV2,
		defaultLogoutURLV2:              config.DefaultLogoutURLV2,
		defaultAccessTokenLifetime:      config.DefaultAccessTokenLifetime,
		defaultIdTokenLifetime:          config.DefaultIdTokenLifetime,
		jwksCacheControlMaxAge:          config.JWKSCacheControlMaxAge,
		fallbackLogger:                  fallbackLogger,
		hasher:                          hasher,
		encAlg:                          authAlg,
		targetEncryptionAlgorithm:       targetEncryptionAlgorithm,
		opCrypto:                        alg,
		assetAPIPrefix:                  assets.AssetAPI(),
		httpClient:                      httpClient,
		registrationEndpoint:            registrationEndpoint(config.CustomEndpoints),
		dynamicClientRegistrationConfig: config.DynamicClientRegistration,
		clientIDMetadataResolver:        clientIDMetadataResolver,
	}
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}

	// We register the routes via op.RegisterServer (instead of op.RegisterLegacyServer)
	// so that the dynamic client registration route can be added through the same
	// WithSetRouter hook as the authorize callback. op.RegisterLegacyServer appends its
	// own WithHTTPMiddleware after the caller options, which would violate chi's
	// "all middlewares before routes" rule once we add a route via WithSetRouter.
	// This mirrors op.RegisterLegacyServer; op.NewIssuerInterceptor reproduces the issuer
	// middleware it would otherwise add.
	server.Handler = op.RegisterServer(server, server.Endpoints(),
		op.WithFallbackLogger(fallbackLogger),
		op.WithHTTPMiddleware(
			middleware.CallDurationHandler,
			middleware.RequestDetailsHandler(),
			middleware.MetricsHandler(metricTypes),
			middleware.TraceHandler(),
			middleware.LogHandler("oidc"),
			middleware.RecoverHandler(writeRecoverError),
			middleware.NoCacheInterceptor().Handler,
			instanceHandler,
			userAgentCookie,
			http_utils.CopyHeadersToContext,
			accessHandler.HandleWithPublicAuthPathPrefixes(publicAuthPathPrefixes(config.CustomEndpoints)),
			middleware.ActivityHandler,
		),
		op.WithHTTPMiddleware(op.NewIssuerInterceptor(server.IssuerFromRequest).Handler),
		op.WithSetRouter(func(r chi.Router) {
			r.HandleFunc(server.Endpoints().Authorization.Relative()+authCallbackPathSuffix, server.authorizeCallbackHandler)
			r.Method(http.MethodPost, server.registrationEndpoint.Relative(), http.HandlerFunc(server.dynamicClientRegistration))
		}),
	)

	return server, nil
}

// authCallbackPathSuffix mirrors the unexported suffix used by op.RegisterLegacyServer to
// register the authorize callback handler under the authorization endpoint.
// Keep in sync with github.com/zitadel/oidc/v3/pkg/op (authCallbackPathSuffix). The
// existing authorization-flow integration tests exercise this route and would fail if the
// library changed the suffix.
const authCallbackPathSuffix = "/callback"

// registrationEndpoint builds the dynamic client registration endpoint, optionally
// overridden through the custom endpoint configuration.
func registrationEndpoint(endpointConfig *EndpointConfig) *op.Endpoint {
	if endpointConfig != nil && endpointConfig.Registration != nil {
		return op.NewEndpointWithURL(endpointConfig.Registration.Path, endpointConfig.Registration.URL)
	}
	return op.NewEndpoint(defaultRegistrationEndpoint)
}

func ContextToIssuer(ctx context.Context) string {
	return http_utils.DomainContext(ctx).Origin()
}

func IssuerFromContext(_ bool) (op.IssuerFromRequest, error) {
	return func(r *http.Request) string {
		return ContextToIssuer(r.Context())
	}, nil
}

func publicAuthPathPrefixes(endpoints *EndpointConfig) []string {
	authURL := op.DefaultEndpoints.Authorization.Relative()
	keysURL := op.DefaultEndpoints.JwksURI.Relative()
	if endpoints == nil {
		return []string{oidc.DiscoveryEndpoint, authURL, keysURL}
	}
	if endpoints.Auth != nil && endpoints.Auth.Path != "" {
		authURL = endpoints.Auth.Path
	}
	if endpoints.Keys != nil && endpoints.Keys.Path != "" {
		keysURL = endpoints.Keys.Path
	}
	return []string{oidc.DiscoveryEndpoint, authURL, keysURL}
}

func createOPConfig(config Config, defaultLogoutRedirectURI string, cryptoKey []byte) (*op.Config, error) {
	opConfig := &op.Config{
		DefaultLogoutRedirectURI: defaultLogoutRedirectURI,
		CodeMethodS256:           config.CodeMethodS256,
		AuthMethodPost:           config.AuthMethodPost,
		AuthMethodPrivateKeyJWT:  config.AuthMethodPrivateKeyJWT,
		GrantTypeRefreshToken:    config.GrantTypeRefreshToken,
		RequestObjectSupported:   config.RequestObjectSupported,
		DeviceAuthorization:      config.DeviceAuth.toOPConfig(),
	}
	if cryptoLength := len(cryptoKey); cryptoLength != 32 {
		return nil, zerrors.ThrowInternalf(nil, "OIDC-D43gf", "crypto key must be 32 bytes, but is %d", cryptoLength)
	}
	copy(opConfig.CryptoKey[:], cryptoKey)
	return opConfig, nil
}

func newStorage(
	config Config,
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	authAlg crypto.AuthAlgorithm,
	es *eventstore.Eventstore,
	contextToIssuer func(context.Context) string,
	federateLogoutCache cache.Cache[federatedlogout.Index, string, *federatedlogout.FederatedLogout],
	clientIDMetadataResolver *clientIDMetadataResolver,
) *OPStorage {
	return &OPStorage{
		repo:                              repo,
		command:                           command,
		query:                             query,
		eventstore:                        es,
		defaultLoginURL:                   fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID),
		defaultLoginURLV2:                 config.DefaultLoginURLV2,
		defaultLogoutURLV2:                config.DefaultLogoutURLV2,
		defaultAccessTokenLifetime:        config.DefaultAccessTokenLifetime,
		defaultIdTokenLifetime:            config.DefaultIdTokenLifetime,
		defaultRefreshTokenIdleExpiration: config.DefaultRefreshTokenIdleExpiration,
		defaultRefreshTokenExpiration:     config.DefaultRefreshTokenExpiration,
		authAlg:                           authAlg,
		assetAPIPrefix:                    assets.AssetAPI(),
		contextToIssuer:                   contextToIssuer,
		federateLogoutCache:               federateLogoutCache,
		clientIDMetadataResolver:          clientIDMetadataResolver,
	}
}

func (o *OPStorage) Health(ctx context.Context) error {
	return o.repo.Health(ctx)
}
