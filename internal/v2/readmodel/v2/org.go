package readmodel

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/org"
	"github.com/zitadel/zitadel/internal/v2/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var _ object = (*Org)(nil)

type Org struct {
	ID   string
	Name string

	PrimaryDomain *projection.OrgPrimaryDomain
	State         *projection.OrgState

	projection.ObjectMetadata
}

func NewOrg(id string) *Org {
	return &Org{
		ID:            id,
		PrimaryDomain: projection.NewOrgPrimaryDomain(id),
		State:         projection.NewOrgState(id),
	}
}

func (o *Org) Reducers() projection.Reducers {
	if o.ObjectMetadata.Reducers != nil {
		return o.ObjectMetadata.Reducers
	}

	o.ObjectMetadata.Reducers = projection.MergeReducers(
		projection.Reducers{
			org.AggregateType: {
				org.AddedType:   o.reduceAdded,
				org.ChangedType: o.reduceChanged,
			},
		},
		o.PrimaryDomain.Reducers(),
		o.State.Reducers(),
	)

	return o.ObjectMetadata.Reducers
}

func (o *Org) reduceAdded(event *v2_es.StorageEvent) error {
	if !o.ObjectMetadata.ShouldReduce(event) {
		return nil
	}

	e, err := org.AddedEventFromStorage(event)
	if err != nil {
		return err
	}

	o.Name = e.Payload.Name
	o.ObjectMetadata.Set(event)
	return nil
}

func (o *Org) reduceChanged(event *v2_es.StorageEvent) error {
	if !o.ObjectMetadata.ShouldReduce(event) {
		return nil
	}

	e, err := org.ChangedEventFromStorage(event)
	if err != nil {
		return err
	}

	o.Name = e.Payload.Name
	o.ObjectMetadata.Set(event)
	return nil
}

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

func (o *Orgs) ByID(ctx context.Context, id string) (*Org, error) {
	org, ok := o.cache.Get(id)
	if !ok {
		if err := o.loadOrg(ctx, id); err != nil {
			return nil, err
		}
		org, ok = o.cache.Get(id)
		if !ok {
			return nil, zerrors.ThrowNotFound(nil, "V2-rUKyI", "Errors.Org.NotFound")
		}
	}
	return org, nil
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
func (o *Orgs) Reducers() projection.Reducers {
	if o.reducers != nil {
		return o.reducers
	}

	o.reducers = projection.OverwriteReduces(NewOrg("").Reducers(), o.reduce)
	o.reducers[org.AggregateType][org.AddedType] = o.reduceAdded

	return o.reducers
}

func (o *Orgs) reduce(event *v2_es.StorageEvent) error {
	org, ok := o.cache.Get(event.Aggregate.ID)
	if !ok {
		return nil
	}
	err := org.ObjectMetadata.Reduce(event, org.Reducers()[event.Aggregate.Type][event.Type])
	if err != nil {
		return err
	}
	return o.cache.Set(org.ID, org)
}

func (o *Orgs) reduceAdded(event *v2_es.StorageEvent) error {
	org, ok := o.cache.Get(event.Aggregate.ID)
	if !ok {
		org = &Org{
			ID:            event.Aggregate.ID,
			PrimaryDomain: projection.NewOrgPrimaryDomain(event.Aggregate.ID),
			State:         projection.NewOrgState(event.Aggregate.ID),
		}
	}
	err := org.ObjectMetadata.Reduce(event, org.Reducers()[event.Aggregate.Type][event.Type])
	if err != nil {
		return err
	}
	return o.cache.Set(org.ID, org)
}

func (o *Orgs) loadOrg(ctx context.Context, id string) error {
	return o.list.es.FilterToReducer(ctx, o.EventstoreV3Query(decimal.Zero), o.list)
}
