package readmodel

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
)

var _ object = (*Org)(nil)

type Org struct{}

type Orgs struct {
	readModel

	cache Cache[string, *Org]
	list  *listReadModel
}

func NewOrgs(ctx context.Context, es *eventstore.Eventstore) *Orgs {
	orgs := &Orgs{
		cache: NewMapCache[string, *Org](),
	}
	orgs.list = newListReadModel(ctx, orgs, es)
	return orgs
}

// EventstoreV3Query implements listManager.
func (o *Orgs) EventstoreV3Query(position decimal.Decimal) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(eventstore.AggregateType(org.AggregateType)).
		Builder()
}

// Name implements listManager.
func (o *Orgs) Name() string {
	return "orgs"
}

// Reducers implements listManager.
func (o *Orgs) Reducers() map[string]map[string]eventstore.ReduceEvent {
	panic("unimplemented")
}
