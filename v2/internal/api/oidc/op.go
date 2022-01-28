package oidc

import (
	"context"
	"database/sql"
	"time"

	"github.com/caos/logging"
	"github.com/caos/oidc/pkg/op"
	"golang.org/x/text/language"

	http_utils "github.com/caos/zitadel/internal/api/http"

	"github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/auth/repository"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type OPHandlerConfig struct {
	OPConfig              *op.Config
	StorageConfig         StorageConfig
	UserAgentCookieConfig *middleware.UserAgentCookieConfig
	Cache                 *middleware.CacheConfig
	Endpoints             *EndpointConfig
}

type StorageConfig struct {
	DefaultLoginURL                   string
	SigningKeyAlgorithm               string
	DefaultAccessTokenLifetime        types.Duration
	DefaultIdTokenLifetime            types.Duration
	DefaultRefreshTokenIdleExpiration types.Duration
	DefaultRefreshTokenExpiration     types.Duration
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

func NewProvider(ctx context.Context, config OPHandlerConfig, command *command.Commands, query *query.Queries, repo repository.Repository, keyConfig systemdefaults.KeyConfig, localDevMode bool, es *eventstore.Eventstore, projections *sql.DB, keyChan <-chan interface{}, assetAPIPrefix string) op.OpenIDProvider {
	cookieHandler, err := middleware.NewUserAgentHandler(config.UserAgentCookieConfig, id.SonyFlakeGenerator, localDevMode)
	logging.Log("OIDC-sd4fd").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Panic("cannot user agent handler")
	tokenKey, err := crypto.LoadKey(keyConfig.EncryptionConfig, keyConfig.EncryptionConfig.EncryptionKeyID)
	logging.Log("OIDC-ADvbv").OnError(err).Panic("cannot load OP crypto key")
	cryptoKey := []byte(tokenKey)
	if len(cryptoKey) != 32 {
		logging.Log("OIDC-Dsfds").Panic("OP crypto key must be exactly 32 bytes")
	}
	copy(config.OPConfig.CryptoKey[:], cryptoKey)
	supportedLanguages, err := getSupportedLanguages()
	logging.Log("OIDC-GBd3t").OnError(err).Panic("cannot get supported languages")
	config.OPConfig.SupportedUILocales = supportedLanguages
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}
	storage, err := newStorage(config.StorageConfig, command, query, repo, keyConfig, es, projections, keyChan, assetAPIPrefix)
	logging.Log("OIDC-Jdg2k").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Panic("cannot create storage")
	options := []op.Option{
		op.WithHttpInterceptors(
			middleware.MetricsHandler(metricTypes),
			middleware.TelemetryHandler(),
			middleware.NoCacheInterceptor,
			cookieHandler,
			http_utils.CopyHeadersToContext,
		),
	}
	options = append(options, customEndpoints(config.Endpoints)...)
	provider, err := op.NewOpenIDProvider(
		ctx,
		config.OPConfig,
		storage,
		options...,
	)
	logging.Log("OIDC-asf13").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Panic("cannot create provider")
	return provider
}

func customEndpoints(endpointConfig *EndpointConfig) []op.Option {
	if endpointConfig == nil {
		return nil
	}
	options := []op.Option{}
	if endpointConfig.Auth != nil {
		options = append(options, op.WithCustomAuthEndpoint(op.NewEndpointWithURL(endpointConfig.Auth.Path, endpointConfig.Auth.URL)))
	}
	if endpointConfig.Auth != nil {
		options = append(options, op.WithCustomTokenEndpoint(op.NewEndpointWithURL(endpointConfig.Token.Path, endpointConfig.Token.URL)))
	}
	if endpointConfig.Auth != nil {
		options = append(options, op.WithCustomIntrospectionEndpoint(op.NewEndpointWithURL(endpointConfig.Introspection.Path, endpointConfig.Introspection.URL)))
	}
	if endpointConfig.Auth != nil {
		options = append(options, op.WithCustomUserinfoEndpoint(op.NewEndpointWithURL(endpointConfig.Userinfo.Path, endpointConfig.Userinfo.URL)))
	}
	if endpointConfig.Auth != nil {
		options = append(options, op.WithCustomRevocationEndpoint(op.NewEndpointWithURL(endpointConfig.Revocation.Path, endpointConfig.Revocation.URL)))
	}
	if endpointConfig.Auth != nil {
		options = append(options, op.WithCustomEndSessionEndpoint(op.NewEndpointWithURL(endpointConfig.EndSession.Path, endpointConfig.EndSession.URL)))
	}
	if endpointConfig.Auth != nil {
		options = append(options, op.WithCustomKeysEndpoint(op.NewEndpointWithURL(endpointConfig.Keys.Path, endpointConfig.Keys.URL)))
	}
	return options
}

func newStorage(config StorageConfig, command *command.Commands, query *query.Queries, repo repository.Repository, keyConfig systemdefaults.KeyConfig, es *eventstore.Eventstore, projections *sql.DB, keyChan <-chan interface{}, assetAPIPrefix string) (*OPStorage, error) {
	encAlg, err := crypto.NewAESCrypto(keyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	return &OPStorage{
		repo:                              repo,
		command:                           command,
		query:                             query,
		eventstore:                        es,
		defaultLoginURL:                   config.DefaultLoginURL,
		signingKeyAlgorithm:               config.SigningKeyAlgorithm,
		defaultAccessTokenLifetime:        config.DefaultAccessTokenLifetime.Duration,
		defaultIdTokenLifetime:            config.DefaultIdTokenLifetime.Duration,
		defaultRefreshTokenIdleExpiration: config.DefaultRefreshTokenIdleExpiration.Duration,
		defaultRefreshTokenExpiration:     config.DefaultRefreshTokenExpiration.Duration,
		encAlg:                            encAlg,
		signingKeyGracefulPeriod:          keyConfig.SigningKeyGracefulPeriod.Duration,
		signingKeyRotationCheck:           keyConfig.SigningKeyRotationCheck.Duration,
		locker:                            crdb.NewLocker(projections, locksTable, signingKey),
		keyChan:                           keyChan,
		assetAPIPrefix:                    assetAPIPrefix,
	}, nil
}

func (o *OPStorage) Health(ctx context.Context) error {
	return o.repo.Health(ctx)
}

func getSupportedLanguages() ([]language.Tag, error) {
	//statikLoginFS, err := fs.NewWithNamespace("login")
	//if err != nil {
	//	return nil, err
	//}
	//return i18n.SupportedLanguages(i18nDir)
	return nil, nil //TODO: handle embed
}
