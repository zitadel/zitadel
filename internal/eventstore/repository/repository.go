package repository

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/models"
)

type Repository interface {
	Health(ctx context.Context) error

	// PushEvents adds all events of the given aggregates to the eventstreams of the aggregates.
	// This call is transaction save. The transaction will be rolled back if one event fails
	PushEvents(ctx context.Context, aggregates ...*models.Aggregate) error
	// Filter returns all events matching the given search query
	Filter(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
}
