package repository

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
)

type Repository interface {
	Health(ctx context.Context) error

	// Filter returns all events matching the given search query
	Filter(ctx context.Context, searchQuery *models.SearchQueryFactory) (events []*models.Event, err error)
	//LatestSequence returns the latests sequence found by the the search query
	LatestSequence(ctx context.Context, queryFactory *models.SearchQueryFactory) (uint64, error)
}
