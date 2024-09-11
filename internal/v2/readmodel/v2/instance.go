package readmodel

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/zitadel/zitadel/internal/eventstore"
	v2_es "github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/v2/instance"
	"github.com/zitadel/zitadel/internal/v2/projection"
)

var _ object = (*Instance)(nil)

type Instance struct {
	*projection.Instance
}

var _ listManager = (*Instances)(nil)

// Instances is the manager for the instance list read model.
type Instances struct {
	readModel

	cache Cache[string, *Instance]
	list  *listReadModel
}

func NewInstances(ctx context.Context, es *eventstore.Eventstore) *Instances {
	instances := &Instances{
		cache: NewMapCache[string, *Instance](),
	}
	instances.list = newListReadModel(ctx, instances, es)
	return instances
}

// Name implements manager.
func (i *Instances) Name() string {
	return "instances"
}

// EventstoreV3Query implements manager.
func (i *Instances) EventstoreV3Query(position decimal.Decimal) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(eventstore.AggregateType(instance.AggregateType)).
		EventTypes().
		Builder()
}

// Reducers implements manager.
func (i *Instances) Reducers() projection.Reducers {
	if i.reducers != nil {
		return i.reducers
	}

	i.reducers = projection.Reducers{
		instance.AggregateType: {
			instance.AddedType:              i.reduceAdded,
			instance.ChangedType:            i.reduce,
			instance.DefaultOrgSetType:      i.reduce,
			instance.ProjectSetType:         i.reduce,
			instance.ConsoleSetType:         i.reduce,
			instance.DefaultLanguageSetType: i.reduce,
			instance.RemovedType:            i.reduce,

			instance.DomainAddedType:      i.reduce,
			instance.DomainVerifiedType:   i.reduce,
			instance.DomainPrimarySetType: i.reduce,
			instance.DomainRemovedType:    i.reduce,
		},
	}

	return i.reducers
}

func (i *Instances) reduceAdded(event *v2_es.StorageEvent) error {
	instance, ok := i.cache.Get(event.Aggregate.ID)
	if !ok {
		instance = &Instance{
			Instance: projection.NewInstanceFromEvent(event),
		}
	}
	err := instance.Reducers()[event.Aggregate.Type][event.Type](event)
	if err != nil {
		return err
	}
	return i.cache.Set(instance.ID, instance)
}

func (i *Instances) reduce(event *v2_es.StorageEvent) error {
	instance, ok := i.cache.Get(event.Aggregate.ID)
	if !ok {
		return nil
	}
	err := instance.Reducers()[event.Aggregate.Type][event.Type](event)
	if err != nil {
		return err
	}
	return i.cache.Set(instance.ID, instance)
}
