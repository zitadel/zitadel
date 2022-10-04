package repository

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type Repository interface {
	Health(ctx context.Context) error

	// Filter returns all events matching the given search query
	Filter(ctx context.Context, searchQuery *models.SearchQueryFactory) (events []*models.Event, err error)
	//LatestCreationDate returns the latest creation date found by the search query
	LatestCreationDate(ctx context.Context, queryFactory *models.SearchQueryFactory) (time.Time, error)
	//InstanceIDs returns the instance ids found by the search query
	InstanceIDs(ctx context.Context, queryFactory *models.SearchQueryFactory) ([]string, error)
}
