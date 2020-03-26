package eventstore

import (
	"context"

	lib "github.com/caos/eventstore-lib"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type App interface {
	Health(ctx context.Context) error
	CreateEvents(ctx context.Context, aggregates ...*models.Aggregate) (err error)
	FilterEvents(ctx context.Context, searchQuery *models.SearchQuery) (events *models.Events, err error)
}

var _ App = (*app)(nil)

type app struct {
	eventstore lib.Eventstore
}
