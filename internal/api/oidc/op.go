package oidc

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/caos/oidc/pkg/op"
	"github.com/rakyll/statik/fs"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/assets"
	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/api/ui/login"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/i18n"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/metrics"
)

const (
	HandlerPrefix = "/oauth/v2"
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
	keyChan                           <-chan interface{}
	currentKey                        query.PrivateKey
	signingKeyRotationCheck           time.Duration
	signingKeyGracefulPeriod          time.Duration
	locker                            crdb.Locker
	assetAPIPrefix                    string
}

func NewProvider(ctx context.Context, config Config, issuer, defaultLogoutRedirectURI string, command *command.Commands, query *query.Queries, repo repository.Repository, keyConfig systemdefaults.KeyConfig, encryptionAlg crypto.EncryptionAlgorithm, cryptoKey []byte, es *eventstore.Eventstore, projections *sql.DB, keyChan <-chan interface{}, userAgentCookie, instanceHandler func(http.Handler) http.Handler) (op.OpenIDProvider, error) {
	opConfig, err := createOPConfig(config, issuer, defaultLogoutRedirectURI, cryptoKey)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "OIDC-EGrqd", "cannot create op config: %w")
	}
	storage := newStorage(config, command, query, repo, keyConfig, encryptionAlg, es, projections, keyChan)
	options, err := createOptions(config, userAgentCookie, instanceHandler)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "OIDC-D3gq1", "cannot create options: %w")
	}
	provider, err := op.NewOpenIDProvider(
		ctx,
		opConfig,
		storage,
		options...,
	)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "OIDC-DAtg3", "cannot create provider")
	}
	return provider, nil
}

func Issuer(domain string, port uint16, externalSecure bool) string {
	return http_utils.BuildHTTP(domain, port, externalSecure) + HandlerPrefix
}

func createOPConfig(config Config, issuer, defaultLogoutRedirectURI string, cryptoKey []byte) (*op.Config, error) {
	supportedLanguages, err := getSupportedLanguages()
	if err != nil {
		return nil, err
	}
	opConfig := &op.Config{
		Issuer:                   issuer,
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

func createOptions(config Config, userAgentCookie, instanceHandler func(http.Handler) http.Handler) ([]op.Option, error) {
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}
	interceptor := op.WithHttpInterceptors(
		middleware.MetricsHandler(metricTypes),
		middleware.TelemetryHandler(),
		middleware.NoCacheInterceptor,
		instanceHandler,
		userAgentCookie,
		http_utils.CopyHeadersToContext,
	)
	endpoints := customEndpoints(config.CustomEndpoints)
	if len(endpoints) == 0 {
		return []op.Option{interceptor}, nil
	}
	return append(endpoints, interceptor), nil
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

func newStorage(config Config, command *command.Commands, query *query.Queries, repo repository.Repository, keyConfig systemdefaults.KeyConfig, encAlg crypto.EncryptionAlgorithm, es *eventstore.Eventstore, projections *sql.DB, keyChan <-chan interface{}) *OPStorage {
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
		signingKeyGracefulPeriod:          keyConfig.SigningKeyGracefulPeriod,
		signingKeyRotationCheck:           keyConfig.SigningKeyRotationCheck,
		locker:                            crdb.NewLocker(projections, locksTable, signingKey),
		keyChan:                           keyChan,
		assetAPIPrefix:                    assets.HandlerPrefix,
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
