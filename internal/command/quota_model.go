package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
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
		unit: unit,
	}
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
		).EventData(map[string]interface{}{"unit": wm.unit})

	return query.Builder()
}

func (wm *quotaWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *quota.AddedEvent:
			wm.AggregateID = e.Aggregate().ID
			wm.active = true
		case *quota.RemovedEvent:
			wm.AggregateID = e.Aggregate().ID
			wm.active = false
		}
	}
	return wm.WriteModel.Reduce()
}
