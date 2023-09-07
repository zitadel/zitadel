package saml

import (
	"fmt"
	"net/http"

	"github.com/zitadel/saml/pkg/provider"

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
)

const (
	HandlerPrefix = "/saml/v2"
)

type Config struct {
	ProviderConfig *provider.Config
}

func NewProvider(
	conf Config,
	externalSecure bool,
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	encAlg crypto.EncryptionAlgorithm,
	certEncAlg crypto.EncryptionAlgorithm,
	es *eventstore.Eventstore,
	projections *database.DB,
	instanceHandler,
	userAgentCookie func(http.Handler) http.Handler,
	accessHandler *middleware.AccessInterceptor,
) (*provider.Provider, error) {
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}

	provStorage, err := newStorage(
		command,
		query,
		repo,
		encAlg,
		certEncAlg,
		es,
		projections,
	)
	if err != nil {
		return nil, err
	}

	options := []provider.Option{
		provider.WithHttpInterceptors(
			middleware.MetricsHandler(metricTypes),
			middleware.TelemetryHandler(),
			middleware.NoCacheInterceptor().Handler,
			instanceHandler,
			userAgentCookie,
			accessHandler.HandleIgnorePathPrefixes(ignoredQuotaLimitEndpoint(conf.ProviderConfig)),
			http_utils.CopyHeadersToContext,
		),
		provider.WithCustomTimeFormat("2006-01-02T15:04:05.999Z"),
	}
	if !externalSecure {
		options = append(options, provider.WithAllowInsecure())
	}

	return provider.NewProvider(
		provStorage,
		HandlerPrefix,
		conf.ProviderConfig,
		options...,
	)
}

func newStorage(
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	encAlg crypto.EncryptionAlgorithm,
	certEncAlg crypto.EncryptionAlgorithm,
	es *eventstore.Eventstore,
	db *database.DB,
) (*Storage, error) {
	return &Storage{
		encAlg:          encAlg,
		certEncAlg:      certEncAlg,
		locker:          crdb.NewLocker(db.DB, locksTable, signingKey),
		eventstore:      es,
		repo:            repo,
		command:         command,
		query:           query,
		defaultLoginURL: fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID),
	}, nil
}

func ignoredQuotaLimitEndpoint(config *provider.Config) []string {
	ignoredEndpoints := make([]string, 3)
	ignoredEndpoints[0] = HandlerPrefix + provider.DefaultMetadataEndpoint
	ignoredEndpoints[1] = HandlerPrefix + provider.DefaultCertificateEndpoint
	ignoredEndpoints[2] = HandlerPrefix + provider.DefaultSingleSignOnEndpoint
	if config.MetadataConfig != nil && config.MetadataConfig.Path != "" {
		ignoredEndpoints[0] = HandlerPrefix + config.MetadataConfig.Path
	}
	if config.IDPConfig == nil || config.IDPConfig.Endpoints == nil {
		return ignoredEndpoints
	}
	if config.IDPConfig.Endpoints.Certificate != nil && config.IDPConfig.Endpoints.Certificate.Relative() != "" {
		ignoredEndpoints[1] = HandlerPrefix + config.IDPConfig.Endpoints.Certificate.Relative()
	}
	if config.IDPConfig.Endpoints.SingleSignOn != nil && config.IDPConfig.Endpoints.SingleSignOn.Relative() != "" {
		ignoredEndpoints[2] = HandlerPrefix + config.IDPConfig.Endpoints.SingleSignOn.Relative()
	}
	return ignoredEndpoints
}
