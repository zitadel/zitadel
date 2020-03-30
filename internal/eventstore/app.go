package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

type App interface {
	Health(ctx context.Context) error
	CreateEvents(ctx context.Context, aggregates ...*models.Aggregate) (err error)
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events []*models.Event, err error)
}

var _ App = (*app)(nil)

type app struct {
	repo repository.Repository
}
