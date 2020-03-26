package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/models"

	lib_models "github.com/caos/eventstore-lib/pkg/models"
)

func (es *app) CreateEvents(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	libAggregates := make([]lib_models.Aggregate, len(aggregates))
	for idx, aggregate := range aggregates {
		libAggregates[idx] = aggregate
	}

	return es.eventstore.PushEvents(ctx, libAggregates...)
}

func (es *app) FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events *models.Events, err error) {
	events = models.InitEvents()
	if err := es.eventstore.Filter(ctx, events, searchQuery); err != nil {
		return nil, err
	}
	return events, nil
}
