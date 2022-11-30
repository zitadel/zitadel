package projection

import "github.com/zitadel/zitadel/internal/eventstore"

type Projection interface {
	Reduce([]eventstore.Event)
	SearchQuery() *eventstore.SearchQueryBuilder
}
