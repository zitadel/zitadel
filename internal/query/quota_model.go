package query

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type quotaReadModel struct {
	eventstore.ReadModel
	unit   quota.Unit
	active bool
	config *quota.AddedEvent
}

// newQuotaReadModel aggregateId is filled by reducing unit matching events
func newQuotaReadModel(instanceId, resourceOwner string, unit quota.Unit) *quotaReadModel {
	return &quotaReadModel{
		ReadModel: eventstore.ReadModel{
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
		unit: unit,
	}
}

func (rm *quotaReadModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(rm.ResourceOwner).
		AllowTimeTravel().
		AddQuery().
		InstanceID(rm.InstanceID).
		AggregateTypes(quota.AggregateType).
		EventTypes(
			quota.AddedEventType,
			quota.RemovedEventType,
		).EventData(map[string]interface{}{"unit": rm.unit})

	return query.Builder()
}

func (rm *quotaReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *quota.AddedEvent:
			rm.AggregateID = e.Aggregate().ID
			rm.active = true
			rm.config = e
		case *quota.RemovedEvent:
			rm.AggregateID = e.Aggregate().ID
			rm.active = false
			rm.config = nil
		}
	}
	return rm.ReadModel.Reduce()
}
