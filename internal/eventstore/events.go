package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/models"
)

func (es *app) CreateEvents(ctx context.Context, aggregates ...*models.Aggregate) (err error) {
	for _, agg := range aggregates {
		if err = agg.Validate(); err != nil {
			return err
		}
	}
	return es.repo.PushEvents(ctx, aggregates...)
}

func (es *app) FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error) {
	if err = searchQuery.Validate(); err != nil {
		return nil, err
	}
	return es.repo.Filter(ctx, searchQuery)
}
