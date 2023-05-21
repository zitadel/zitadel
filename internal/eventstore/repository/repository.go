package repository

import (
	"context"
	"time"
)

// Repository pushes and filters events
type Repository interface {
	//Health checks if the connection to the storage is available
	Health(ctx context.Context) error
	// Filter returns all events matching the given search query
	Filter(ctx context.Context, searchQuery *SearchQuery) (events []*Event, err error)
	//LatestSequence returns the latest sequence found by the search query
	LatestSequence(ctx context.Context, queryFactory *SearchQuery) (time.Time, error)
	//InstanceIDs returns the instance ids found by the search query
	InstanceIDs(ctx context.Context, queryFactory *SearchQuery) ([]string, error)
	//CreateInstance creates a new sequence for the given instance
	CreateInstance(ctx context.Context, instanceID string) error
}
