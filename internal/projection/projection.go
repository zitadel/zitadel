package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type Projection interface {
	Reduce([]eventstore.Event)
	SearchQuery(context.Context) *eventstore.SearchQueryBuilder
}
