package oidc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rakyll/statik/fs"
	"github.com/zitadel/oidc/v2/pkg/op"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/assets"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
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
	UserAgentCookieConfig             *middleware.UserAgentCookieConfig
	Cache                             *middleware.CacheConfig
	CustomEndpoints                   *EndpointConfig
}

type EndpointConfig struct {
	Auth          *Endpoint
	Token         *Endpoint
	Introspection *Endpoint
	Userinfo      *Endpoint
	Revocation    *Endpoint
	EndSession    *Endpoint
	Keys          *Endpoint
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
	defaultAccessTokenLifetime        time.Duration
	defaultIdTokenLifetime            time.Duration
	signingKeyAlgorithm               string
	defaultRefreshTokenIdleExpiration time.Duration
	defaultRefreshTokenExpiration     time.Duration
	encAlg                            crypto.EncryptionAlgorithm
	locker                            crdb.Locker
	assetAPIPrefix                    func(ctx context.Context) string
}

func NewProvider(ctx context.Context, config Config, defaultLogoutRedirectURI string, externalSecure bool, command *command.Commands, query *query.Queries, repo repository.Repository, encryptionAlg crypto.EncryptionAlgorithm, cryptoKey []byte, es *eventstore.Eventstore, projections *database.DB, userAgentCookie, instanceHandler, accessHandler func(http.Handler) http.Handler) (op.OpenIDProvider, error) {
	opConfig, err := createOPConfig(config, defaultLogoutRedirectURI, cryptoKey)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "OIDC-EGrqd", "cannot create op config: %w")
	}
	storage := newStorage(config, command, query, repo, encryptionAlg, es, projections, externalSecure)
	options, err := createOptions(config, externalSecure, userAgentCookie, instanceHandler, accessHandler)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "OIDC-D3gq1", "cannot create options: %w")
	}
	provider, err := op.NewDynamicOpenIDProvider(
		ctx,
		"",
		opConfig,
		storage,
		options...,
	)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "OIDC-DAtg3", "cannot create provider")
	}
	return provider, nil
}

func createOPConfig(config Config, defaultLogoutRedirectURI string, cryptoKey []byte) (*op.Config, error) {
	supportedLanguages, err := getSupportedLanguages()
	if err != nil {
		return nil, err
	}
	opConfig := &op.Config{
		DefaultLogoutRedirectURI: defaultLogoutRedirectURI,
		CodeMethodS256:           config.CodeMethodS256,
		AuthMethodPost:           config.AuthMethodPost,
		AuthMethodPrivateKeyJWT:  config.AuthMethodPrivateKeyJWT,
		GrantTypeRefreshToken:    config.GrantTypeRefreshToken,
		RequestObjectSupported:   config.RequestObjectSupported,
		SupportedUILocales:       supportedLanguages,
	}
	if cryptoLength := len(cryptoKey); cryptoLength != 32 {
		return nil, caos_errs.ThrowInternalf(nil, "OIDC-D43gf", "crypto key must be 32 bytes, but is %d", cryptoLength)
	}
	copy(opConfig.CryptoKey[:], cryptoKey)
	return opConfig, nil
}

func createOptions(config Config, externalSecure bool, userAgentCookie, instanceHandler, accessHandler func(http.Handler) http.Handler) ([]op.Option, error) {
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}
	options := []op.Option{
		op.WithHttpInterceptors(
			middleware.MetricsHandler(metricTypes),
			middleware.TelemetryHandler(),
			middleware.NoCacheInterceptor().Handler,
			instanceHandler,
			userAgentCookie,
			http_utils.CopyHeadersToContext,
			accessHandler,
		),
	}
	if !externalSecure {
		options = append(options, op.WithAllowInsecure())
	}
	endpoints := customEndpoints(config.CustomEndpoints)
	if len(endpoints) != 0 {
		options = append(options, endpoints...)
	}
	return options, nil
}

func customEndpoints(endpointConfig *EndpointConfig) []op.Option {
	if endpointConfig == nil {
		return nil
	}
	options := []op.Option{}
	if endpointConfig.Auth != nil {
		options = append(options, op.WithCustomAuthEndpoint(op.NewEndpointWithURL(endpointConfig.Auth.Path, endpointConfig.Auth.URL)))
	}
	if endpointConfig.Token != nil {
		options = append(options, op.WithCustomTokenEndpoint(op.NewEndpointWithURL(endpointConfig.Token.Path, endpointConfig.Token.URL)))
	}
	if endpointConfig.Introspection != nil {
		options = append(options, op.WithCustomIntrospectionEndpoint(op.NewEndpointWithURL(endpointConfig.Introspection.Path, endpointConfig.Introspection.URL)))
	}
	if endpointConfig.Userinfo != nil {
		options = append(options, op.WithCustomUserinfoEndpoint(op.NewEndpointWithURL(endpointConfig.Userinfo.Path, endpointConfig.Userinfo.URL)))
	}
	if endpointConfig.Revocation != nil {
		options = append(options, op.WithCustomRevocationEndpoint(op.NewEndpointWithURL(endpointConfig.Revocation.Path, endpointConfig.Revocation.URL)))
	}
	if endpointConfig.EndSession != nil {
		options = append(options, op.WithCustomEndSessionEndpoint(op.NewEndpointWithURL(endpointConfig.EndSession.Path, endpointConfig.EndSession.URL)))
	}
	if endpointConfig.Keys != nil {
		options = append(options, op.WithCustomKeysEndpoint(op.NewEndpointWithURL(endpointConfig.Keys.Path, endpointConfig.Keys.URL)))
	}
	return options
}

func newStorage(config Config, command *command.Commands, query *query.Queries, repo repository.Repository, encAlg crypto.EncryptionAlgorithm, es *eventstore.Eventstore, db *database.DB, externalSecure bool) *OPStorage {
	return &OPStorage{
		repo:                              repo,
		command:                           command,
		query:                             query,
		eventstore:                        es,
		defaultLoginURL:                   fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID),
		signingKeyAlgorithm:               config.SigningKeyAlgorithm,
		defaultAccessTokenLifetime:        config.DefaultAccessTokenLifetime,
		defaultIdTokenLifetime:            config.DefaultIdTokenLifetime,
		defaultRefreshTokenIdleExpiration: config.DefaultRefreshTokenIdleExpiration,
		defaultRefreshTokenExpiration:     config.DefaultRefreshTokenExpiration,
		encAlg:                            encAlg,
		locker:                            crdb.NewLocker(db.DB, locksTable, signingKey),
		assetAPIPrefix:                    assets.AssetAPI(externalSecure),
	}
}

func (o *OPStorage) Health(ctx context.Context) error {
	return o.repo.Health(ctx)
}

func getSupportedLanguages() ([]language.Tag, error) {
	statikLoginFS, err := fs.NewWithNamespace("login")
	if err != nil {
		return nil, err
	}
	return i18n.SupportedLanguages(statikLoginFS)
}
