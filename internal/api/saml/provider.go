package saml

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zitadel/saml/pkg/provider"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/auth/repository"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
	"net/http"
)

const (
	HandlerPrefix = "/saml"
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
	es *eventstore.Eventstore,
	projections *sql.DB,
	instanceHandler,
	userAgentCookie func(http.Handler) http.Handler,
) (*provider.Provider, error) {
	metricTypes := []metrics.MetricType{metrics.MetricTypeRequestCount, metrics.MetricTypeStatusCode, metrics.MetricTypeTotalCount}

	provStorage, err := newStorage(
		command,
		query,
		repo,
		encAlg,
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
			middleware.NoCacheInterceptor,
			instanceHandler,
			userAgentCookie,
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
	es *eventstore.Eventstore,
	projections *sql.DB,
) (*Storage, error) {
	return &Storage{
		encAlg:          encAlg,
		locker:          crdb.NewLocker(projections, locksTable, signingKey),
		eventstore:      es,
		repo:            repo,
		command:         command,
		query:           query,
		defaultLoginURL: fmt.Sprintf("%s%s?%s=", login.HandlerPrefix, login.EndpointLogin, login.QueryAuthRequestID),
	}, nil
}
