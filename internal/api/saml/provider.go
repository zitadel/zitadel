package saml

import (
	"context"
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
	"github.com/zitadel/saml/pkg/provider"
)

const (
	locksTable = "projections.locks"
	signingKey = "signing_key"
)

type ProviderHandlerConfig struct {
	StorageConfig         *StorageConfig `yaml:"StorageConfig"`
	UserAgentCookieConfig *middleware.UserAgentCookieConfig

	ProviderConfig *provider.Config
}

func NewProvider(
	ctx context.Context,
	conf ProviderHandlerConfig,
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	keyConfig systemdefaults.KeyConfig,
	es *eventstore.Eventstore,
	projections types.SQL,
	certChan <-chan interface{},
	localDevMode bool,
) (*provider.Provider, error) {
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}
	cookieHandler, err := middleware.NewUserAgentHandler(conf.UserAgentCookieConfig, id.SonyFlakeGenerator, localDevMode)
	if err != nil {
		return nil, err
	}

	storage, err := newStorage(
		conf.StorageConfig,
		command,
		query,
		repo,
		keyConfig,
		es,
		projections,
		certChan,
	)
	if err != nil {
		return nil, err
	}

	return provider.NewProvider(
		ctx,
		storage,
		conf.ProviderConfig,
		provider.WithHttpInterceptors(
			middleware.MetricsHandler(metricTypes),
			middleware.TelemetryHandler(),
			middleware.NoCacheInterceptor,
			cookieHandler,
			http_utils.CopyHeadersToContext,
		),
	)
}

func newStorage(
	conf *StorageConfig,
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	keyConfig systemdefaults.KeyConfig,
	es *eventstore.Eventstore,
	projections types.SQL,
	certChan <-chan interface{},
) (*Storage, error) {
	encAlg, err := crypto.NewAESCrypto(keyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	sqlClient, err := projections.Start()
	if err != nil {
		return nil, err
	}

	return &Storage{
		certChan:                  certChan,
		certificateAlgorithm:      conf.CertificateAlgorithm,
		encAlg:                    encAlg,
		locker:                    crdb.NewLocker(sqlClient, locksTable, signingKey),
		eventstore:                es,
		certificateRotationCheck:  keyConfig.CertificateRotationCheck.Duration,
		certificateGracefulPeriod: keyConfig.CertificateGracefulPeriod.Duration,
		repo:                      repo,
		command:                   command,
		query:                     query,
		defaultLoginURL:           conf.DefaultLoginURL,
	}, nil
}
