package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/internal/repository"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_v2 "github.com/caos/zitadel/internal/eventstore/v2"
)

type Eventstore interface {
	AggregateCreator() *models.AggregateCreator
	Health(ctx context.Context) error
	PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) error
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
	LatestSequence(ctx context.Context, searchQuery *models.SearchQueryFactory) (uint64, error)
	V2() *es_v2.Eventstore
	Subscribe(aggregates ...models.AggregateType) *Subscription
}

var _ Eventstore = (*eventstore)(nil)

type eventstore struct {
	repo             repository.Repository
	aggregateCreator *models.AggregateCreator

	esV2 *es_v2.Eventstore
}

func (es *eventstore) AggregateCreator() *models.AggregateCreator {
	return es.aggregateCreator
}

func (es *eventstore) PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	for _, aggregate := range aggregates {
		if len(aggregate.Events) == 0 {
			return errors.ThrowInvalidArgument(nil, "EVENT-cNhIj", "no events in aggregate")
		}
		for _, event := range aggregate.Events {
			if err = event.Validate(); err != nil {
				return errors.ThrowInvalidArgument(err, "EVENT-tzIhl", "validate event failed")
			}
		}
	}
	err = es.repo.PushAggregates(ctx, aggregates...)
	if err != nil {
		return err
	}

	go notify(aggregates)
	return nil
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

func (es *eventstore) V2() *es_v2.Eventstore {
	return es.esV2
}
