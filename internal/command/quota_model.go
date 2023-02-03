package command

import (
	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type quotaWriteModel struct {
	eventstore.WriteModel
	unit   quota.Unit
	active bool
}

// newQuotaWriteModel aggregateId is filled by reducing unit matching events
func newQuotaWriteModel(instanceId, resourceOwner string, unit quota.Unit) *quotaWriteModel {
	return &quotaWriteModel{
		WriteModel: eventstore.WriteModel{
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
		unit:   unit,
		active: false,
	}
}

func newQuotaAggregate(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, quota.AggregateType, quota.AggregateVersion)
}

func (wm *quotaWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		InstanceID(wm.InstanceID).
		AggregateTypes(quota.AggregateType).
		EventTypes(
			quota.AddedEventType,
			quota.RemovedEventType,
			quota.NotifiedEventType,
		).EventData(map[string]interface{}{"unit": wm.unit})

	return query.Builder()
}

func (wm *quotaWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *quota.AddedEvent:
			wm.active = true
			wm.AggregateID = e.Aggregate().ID
		case *quota.RemovedEvent:
			wm.active = false
			wm.AggregateID = e.Aggregate().ID
		case *quota.NotifiedEvent:
			wm.AggregateID = e.Aggregate().ID
		}
	}
	return wm.WriteModel.Reduce()
}
