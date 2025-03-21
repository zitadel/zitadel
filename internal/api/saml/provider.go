package saml

import (
	"context"
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
	ProviderConfig    *provider.Config
	DefaultLoginURLV2 string
}

type Provider struct {
	*provider.Provider
	command *command.Commands
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
) (*Provider, error) {
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}

	provStorage, err := newStorage(
		command,
		query,
		repo,
		encAlg,
		certEncAlg,
		es,
		projections,
		fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID),
		conf.DefaultLoginURLV2,
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
			accessHandler.HandleWithPublicAuthPathPrefixes(publicAuthPathPrefixes(conf.ProviderConfig)),
			http_utils.CopyHeadersToContext,
			middleware.ActivityHandler,
		),
		provider.WithCustomTimeFormat("2006-01-02T15:04:05.999Z"),
	}
	if !externalSecure {
		options = append(options, provider.WithAllowInsecure())
	}

	p, err := provider.NewProvider(
		provStorage,
		IssuerFromContext,
		conf.ProviderConfig,
		options...,
	)
	if err != nil {
		return nil, err
	}
	return &Provider{
		p,
		command,
	}, nil
}

func ContextToIssuer(ctx context.Context) string {
	return http_utils.DomainContext(ctx).Origin() + HandlerPrefix
}

func IssuerFromContext(_ bool) (provider.IssuerFromRequest, error) {
	return func(r *http.Request) string {
		return ContextToIssuer(r.Context())
	}, nil
}

func newStorage(
	command *command.Commands,
	query *query.Queries,
	repo repository.Repository,
	encAlg crypto.EncryptionAlgorithm,
	certEncAlg crypto.EncryptionAlgorithm,
	es *eventstore.Eventstore,
	db *database.DB,
	defaultLoginURL string,
	defaultLoginURLV2 string,
) (*Storage, error) {
	return &Storage{
		encAlg:            encAlg,
		certEncAlg:        certEncAlg,
		locker:            crdb.NewLocker(db.DB, locksTable, signingKey),
		eventstore:        es,
		repo:              repo,
		command:           command,
		query:             query,
		defaultLoginURL:   defaultLoginURL,
		defaultLoginURLv2: defaultLoginURLV2,
	}, nil
}

func publicAuthPathPrefixes(config *provider.Config) []string {
	metadataEndpoint := HandlerPrefix + provider.DefaultMetadataEndpoint
	certificateEndpoint := HandlerPrefix + provider.DefaultCertificateEndpoint
	ssoEndpoint := HandlerPrefix + provider.DefaultSingleSignOnEndpoint
	if config.MetadataConfig != nil && config.MetadataConfig.Path != "" {
		metadataEndpoint = HandlerPrefix + config.MetadataConfig.Path
	}
	if config.IDPConfig == nil || config.IDPConfig.Endpoints == nil {
		return []string{metadataEndpoint, certificateEndpoint, ssoEndpoint}
	}
	if config.IDPConfig.Endpoints.Certificate != nil && config.IDPConfig.Endpoints.Certificate.Relative() != "" {
		certificateEndpoint = HandlerPrefix + config.IDPConfig.Endpoints.Certificate.Relative()
	}
	if config.IDPConfig.Endpoints.SingleSignOn != nil && config.IDPConfig.Endpoints.SingleSignOn.Relative() != "" {
		ssoEndpoint = HandlerPrefix + config.IDPConfig.Endpoints.SingleSignOn.Relative()
	}
	return []string{metadataEndpoint, certificateEndpoint, ssoEndpoint}
}
