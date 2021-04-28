package repository

import (
	"context"
)

//Repository pushes and filters events
type Repository interface {
	//Health checks if the connection to the storage is available
	Health(ctx context.Context) error
	// PushEvents adds all events of the given aggregates to the eventstreams of the aggregates.
	// if assets are pushed, they will be added to the assets table only referenced by id in eventstore (history)
	// if unique constraints are pushed, they will be added to the unique table for checking unique constraint violations
	// This call is transaction save. The transaction will be rolled back if one event fails
	Push(ctx context.Context, events []*Event, uniqueConstraints ...*UniqueConstraint) error
	// Filter returns all events matching the given search query
	Filter(ctx context.Context, searchQuery *SearchQuery) (events []*Event, err error)
	//LatestSequence returns the latests sequence found by the the search query
	LatestSequence(ctx context.Context, queryFactory *SearchQuery) (uint64, error)
}
