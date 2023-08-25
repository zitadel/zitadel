package v1

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/v1/internal/repository"
	z_sql "github.com/zitadel/zitadel/internal/eventstore/v1/internal/repository/sql"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

type Eventstore interface {
	Health(ctx context.Context) error
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
	Subscribe(aggregates ...models.AggregateType) *Subscription
	InstanceIDs(ctx context.Context, searchQuery *models.SearchQuery) ([]string, error)
}

var _ Eventstore = (*eventstore)(nil)

type eventstore struct {
	repo repository.Repository
}

func Start(db *database.DB, allowOrderByCreationDate bool) (Eventstore, error) {
	return &eventstore{
		repo: z_sql.Start(db, allowOrderByCreationDate),
	}, nil
}

func (es *eventstore) FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (_ []*models.Event, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err := searchQuery.Validate(); err != nil {
		return nil, err
	}
	return es.repo.Filter(ctx, models.FactoryFromSearchQuery(searchQuery))
}

func (es *eventstore) Health(ctx context.Context) error {
	return es.repo.Health(ctx)
}

func (es *eventstore) InstanceIDs(ctx context.Context, searchQuery *models.SearchQuery) ([]string, error) {
	if err := searchQuery.Validate(); err != nil {
		return nil, err
	}
	return es.repo.InstanceIDs(ctx, models.FactoryFromSearchQuery(searchQuery))
}
