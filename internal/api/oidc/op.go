package oidc

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"

	"github.com/zitadel/zitadel/internal/api/assets"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	CodeMethodS256                    bool
	AuthMethodPost                    bool
	AuthMethodPrivateKeyJWT           bool
	GrantTypeRefreshToken             bool
	RequestObjectSupported            bool
	SigningKeyAlgorithm               string
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
	signingKeyAlgorithm               string
	defaultRefreshTokenIdleExpiration time.Duration
	defaultRefreshTokenExpiration     time.Duration
	encAlg                            crypto.EncryptionAlgorithm
	locker                            crdb.Locker
	assetAPIPrefix                    func(ctx context.Context) string
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
		supportedSigningAlgs(ctx)...,
	))
}

// AccessTokenVerifier configures a Verifier and supported signing algorithms based on the Web Key feature in the context.
func (o *Provider) AccessTokenVerifier(ctx context.Context) *op.AccessTokenVerifier {
	return op.NewAccessTokenVerifier(op.IssuerFromContext(ctx), o.accessTokenKeySet, op.WithSupportedAccessTokenSigningAlgorithms(
		supportedSigningAlgs(ctx)...,
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
	encryptionAlg crypto.EncryptionAlgorithm,
	cryptoKey []byte,
	es *eventstore.Eventstore,
	projections *database.DB,
	userAgentCookie, instanceHandler func(http.Handler) http.Handler,
	accessHandler *middleware.AccessInterceptor,
	fallbackLogger *slog.Logger,
	hashConfig crypto.HashConfig,
) (*Server, error) {
	opConfig, err := createOPConfig(config, defaultLogoutRedirectURI, cryptoKey)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-EGrqd", "cannot create op config: %w")
	}
	storage := newStorage(config, command, query, repo, encryptionAlg, es, projections)
	keyCache := newPublicKeyCache(ctx, config.PublicKeyCacheMaxAge, queryKeyFunc(query))
	accessTokenKeySet := newOidcKeySet(keyCache, withKeyExpiryCheck(true))
	idTokenHintKeySet := newOidcKeySet(keyCache)

	var options []op.Option
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
		repo:                       repo,
		query:                      query,
		command:                    command,
		accessTokenKeySet:          accessTokenKeySet,
		idTokenHintKeySet:          idTokenHintKeySet,
		defaultLoginURL:            fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID),
		defaultLoginURLV2:          config.DefaultLoginURLV2,
		defaultLogoutURLV2:         config.DefaultLogoutURLV2,
		defaultAccessTokenLifetime: config.DefaultAccessTokenLifetime,
		defaultIdTokenLifetime:     config.DefaultIdTokenLifetime,
		jwksCacheControlMaxAge:     config.JWKSCacheControlMaxAge,
		fallbackLogger:             fallbackLogger,
		hasher:                     hasher,
		signingKeyAlgorithm:        config.SigningKeyAlgorithm,
		encAlg:                     encryptionAlg,
		opCrypto:                   op.NewAESCrypto(opConfig.CryptoKey),
		assetAPIPrefix:             assets.AssetAPI(),
	}
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}
	server.Handler = op.RegisterLegacyServer(server,
		server.authorizeCallbackHandler,
		op.WithFallbackLogger(fallbackLogger),
		op.WithHTTPMiddleware(
			middleware.MetricsHandler(metricTypes),
			middleware.TelemetryHandler(),
			middleware.NoCacheInterceptor().Handler,
			instanceHandler,
			userAgentCookie,
			http_utils.CopyHeadersToContext,
			accessHandler.HandleWithPublicAuthPathPrefixes(publicAuthPathPrefixes(config.CustomEndpoints)),
			middleware.ActivityHandler,
		))

	return server, nil
}

func IssuerFromContext(_ bool) (op.IssuerFromRequest, error) {
	return func(r *http.Request) string {
		return http_utils.DomainContext(r.Context()).Origin()
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

func newStorage(config Config, command *command.Commands, query *query.Queries, repo repository.Repository, encAlg crypto.EncryptionAlgorithm, es *eventstore.Eventstore, db *database.DB) *OPStorage {
	return &OPStorage{
		repo:                              repo,
		command:                           command,
		query:                             query,
		eventstore:                        es,
		defaultLoginURL:                   fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID),
		defaultLoginURLV2:                 config.DefaultLoginURLV2,
		defaultLogoutURLV2:                config.DefaultLogoutURLV2,
		signingKeyAlgorithm:               config.SigningKeyAlgorithm,
		defaultAccessTokenLifetime:        config.DefaultAccessTokenLifetime,
		defaultIdTokenLifetime:            config.DefaultIdTokenLifetime,
		defaultRefreshTokenIdleExpiration: config.DefaultRefreshTokenIdleExpiration,
		defaultRefreshTokenExpiration:     config.DefaultRefreshTokenExpiration,
		encAlg:                            encAlg,
		locker:                            crdb.NewLocker(db.DB, locksTable, signingKey),
		assetAPIPrefix:                    assets.AssetAPI(),
	}
}

func (o *OPStorage) Health(ctx context.Context) error {
	return o.repo.Health(ctx)
}
