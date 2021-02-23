package repository

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Repository interface {
	Health(ctx context.Context) error

	// PushEvents adds all events of the given aggregates to the eventstreams of the aggregates.
	// This call is transaction save. The transaction will be rolled back if one event fails
	PushAggregates(ctx context.Context, aggregates ...*models.Aggregate) error
	// Filter returns all events matching the given search query
	Filter(ctx context.Context, searchQuery *models.SearchQueryFactory) (events []*models.Event, err error)
	//LatestSequence returns the latests sequence found by the the search query
	LatestSequence(ctx context.Context, queryFactory *models.SearchQueryFactory) (uint64, error)
}
