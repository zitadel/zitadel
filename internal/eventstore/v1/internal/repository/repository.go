package repository

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Repository interface {
	Health(ctx context.Context) error

	// Filter returns all events matching the given search query
	Filter(ctx context.Context, searchQuery *models.SearchQueryFactory) (events []*models.Event, err error)
	//LatestSequence returns the latest sequence found by the search query
	LatestSequence(ctx context.Context, queryFactory *models.SearchQueryFactory) (uint64, error)
	//InstanceIDs returns the instance ids found by the search query
	InstanceIDs(ctx context.Context, queryFactory *models.SearchQueryFactory) ([]string, error)
}
