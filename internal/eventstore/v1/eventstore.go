package v1

import (
	"context"
	"database/sql"

	"github.com/caos/zitadel/internal/eventstore/v1/internal/repository"
	z_sql "github.com/caos/zitadel/internal/eventstore/v1/internal/repository/sql"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Eventstore interface {
	Health(ctx context.Context) error
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
	Subscribe(aggregates ...models.AggregateType) *Subscription
}

var _ Eventstore = (*eventstore)(nil)

type eventstore struct {
	repo repository.Repository
}

func Start(db *sql.DB) (Eventstore, error) {
	return &eventstore{
		repo: z_sql.Start(db),
	}, nil
}

func (es *eventstore) FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) ([]*models.Event, error) {
	if err := searchQuery.Validate(); err != nil {
		return nil, err
	}
	return es.repo.Filter(ctx, models.FactoryFromSearchQuery(searchQuery))
}

func (es *eventstore) Health(ctx context.Context) error {
	return es.repo.Health(ctx)
}
