package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

// resourceOwnerModel can be used to retrieve the resourceOwner of an aggregate
// by checking the first event it.
type resourceOwnerModel struct {
	instanceID    string
	aggregateType eventstore.AggregateType
	aggregateID   string

	resourceOwner string
}

func NewResourceOwnerModel(instanceID string, aggregateType eventstore.AggregateType, aggregateID string) *resourceOwnerModel {
	return &resourceOwnerModel{
		instanceID:    instanceID,
		aggregateType: aggregateType,
		aggregateID:   aggregateID,
	}
}

func (r *resourceOwnerModel) Reduce() error {
	return nil
}
func (r *resourceOwnerModel) AppendEvents(events ...eventstore.Event) {
	if len(events) == 1 {
		r.resourceOwner = events[0].Aggregate().ResourceOwner
	}
}
func (r *resourceOwnerModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		InstanceID(r.instanceID).
		AddQuery().
		AggregateTypes(r.aggregateType).
		AggregateIDs(r.aggregateID).
		Builder().
		Limit(1)
}
