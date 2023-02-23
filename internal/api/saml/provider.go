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
	ProviderConfig *provider.Config
}

func NewProvider(
	ctx context.Context,
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
	userAgentCookie,
	accessHandler func(http.Handler) http.Handler,
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
			accessHandler,
			http_utils.CopyHeadersToContext,
		),
	}
	if !externalSecure {
		options = append(options, provider.WithAllowInsecure())
	}

	return provider.NewProvider(
		ctx,
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
