package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/internal/repository"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Eventstore interface {
	AggregateCreator() *models.AggregateCreator
	Health(ctx context.Context) error
	PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) error
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
}

var _ Eventstore = (*eventstore)(nil)

type eventstore struct {
	repo             repository.Repository
	aggregateCreator *models.AggregateCreator
}

func (es *eventstore) AggregateCreator() *models.AggregateCreator {
	return es.aggregateCreator
}

func (es *eventstore) PushEvents(ctx context.Context, aggregates ...[]*models.Event) (err error) {
	if len(aggregates) == 0 {
		return errors.ThrowInvalidArgument(nil, "EVENT-JCifQ", "no aggregates provided")
	}

	for _, events := range aggregates {
		if len(events) == 0 {
			return errors.ThrowInvalidArgument(nil, "EVENT-eci1Z", "no events provided")
		}
		for _, event := range events {
			if err := event.Validate(); err != nil {
				return errors.ThrowInvalidArgument(err, "EVENT-rsa2U", "event invalid")
			}
		}
	}

	return es.repo.PushEvents(ctx, aggregates)
}

func (es *eventstore) PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	for _, aggregate := range aggregates {
		if err = aggregate.Validate(); err != nil {
			return err
		}
	}
	err = es.repo.PushAggregates(ctx, aggregates...)
	if err != nil {
		return err
	}

	for _, aggregate := range aggregates {
		if aggregate.Appender != nil {
			aggregate.Appender(aggregate.Events...)
		}
	}

	return nil
}

func (es *eventstore) FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) ([]*models.Event, error) {
	if err := searchQuery.Validate(); err != nil {
		return nil, err
	}
	return es.repo.Filter(ctx, searchQuery)
}
