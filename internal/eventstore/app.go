package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/internal/repository"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type App interface {
	Health(ctx context.Context) error
	PushAggregates(ctx context.Context, aggregates ...*Aggregate) ([]*models.Aggregate, error)
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
}

var _ App = (*app)(nil)

type app struct {
	repo repository.Repository
}

func (es *app) PushEvents(ctx context.Context, aggregates ...*Aggregate) (err error) {
	return errors.ThrowUnimplemented(nil, "EVENT-fLtHG", "needs improvement use PushAggregates instead")
	aggs := make([][]*models.Event, len(aggregates))
	for aggIdx, aggregate := range aggregates {
		aggs[aggIdx] = make([]*models.Event, len(aggregate.events))
		for eventIdx, event := range aggregate.events {
			aggs[aggIdx][eventIdx] = &models.Event{
				AggregateID:      event.aggregateID,
				AggregateType:    event.aggregateType,
				AggregateVersion: event.aggregateVersion,
				CreationDate:     event.creationDate,
				Data:             event.data,
				ModifierService:  event.modifierService,
				ModifierTenant:   event.modifierTenant,
				ModifierUser:     event.modifierUser,
				PreviousSequence: event.previousSequence,
				ResourceOwner:    event.resourceOwner,
				Typ:              event.typ,
			}
		}
	}

	return es.repo.PushEvents(ctx, aggs)
}

func (es *app) PushAggregates(ctx context.Context, aggregates ...*Aggregate) (aggs []*models.Aggregate, err error) {
	aggs = make([]*models.Aggregate, len(aggregates))
	for aggIdx, aggregate := range aggregates {
		agg := &models.Aggregate{
			ID:             aggregate.id,
			LatestSequence: aggregate.latestSequence,
			Typ:            aggregate.typ,
			Version:        aggregate.version,
			Events:         make([]*models.Event, len(aggregate.events)),
		}
		for eventIdx, event := range aggregate.events {
			agg.Events[eventIdx] = &models.Event{
				CreationDate:     event.creationDate,
				Data:             event.data,
				ModifierService:  event.modifierService,
				ModifierTenant:   event.modifierTenant,
				ModifierUser:     event.modifierUser,
				PreviousSequence: event.previousSequence,
				ResourceOwner:    event.resourceOwner,
				Typ:              event.typ,
			}
		}
		aggs[aggIdx] = agg
	}
	return aggs, es.repo.PushAggregates(ctx, aggs...)
}

func (es *app) FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error) {
	if err = searchQuery.Validate(); err != nil {
		return nil, err
	}
	return es.repo.Filter(ctx, searchQuery)
}
