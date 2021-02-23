package v1

import (
	"context"
	eventstore2 "github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/v1/internal/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Eventstore interface {
	Health(ctx context.Context) error
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
	LatestSequence(ctx context.Context, searchQuery *models.SearchQueryFactory) (uint64, error)
	V2() *eventstore2.Eventstore
	Subscribe(aggregates ...models.AggregateType) *Subscription
}

var _ Eventstore = (*eventstore)(nil)

type eventstore struct {
	repo             repository.Repository
	aggregateCreator *models.AggregateCreator

	esV2 *eventstore2.Eventstore
}

func (es *eventstore) FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) ([]*models.Event, error) {
	if err := searchQuery.Validate(); err != nil {
		return nil, err
	}
	return es.repo.Filter(ctx, models.FactoryFromSearchQuery(searchQuery))
}

func (es *eventstore) LatestSequence(ctx context.Context, queryFactory *models.SearchQueryFactory) (uint64, error) {
	sequenceFactory := *queryFactory
	sequenceFactory = *(&sequenceFactory).Columns(models.Columns_Max_Sequence)
	sequenceFactory = *(&sequenceFactory).SequenceGreater(0)
	return es.repo.LatestSequence(ctx, &sequenceFactory)
}

func (es *eventstore) Health(ctx context.Context) error {
	return es.repo.Health(ctx)
}

func (es *eventstore) V2() *eventstore2.Eventstore {
	return es.esV2
}
