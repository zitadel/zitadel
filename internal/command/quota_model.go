package command

import (
	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type quotaWriteModel struct {
	eventstore.WriteModel
	unit   quota.Unit
	exists bool
}

func newQuotaWriteModel(aggregateId, instanceId, resourceOwner string) *quotaWriteModel {
	return &quotaWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   aggregateId,
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
	}
}

func newQuotaAggregate(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, quota.AggregateType, quota.AggregateVersion)
}

func (wm *quotaWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		InstanceID(wm.InstanceID).
		AggregateIDs(wm.AggregateID).
		AggregateTypes(quota.AggregateType).
		EventTypes(
			quota.AddedEventType,
			quota.RemovedEventType,
			quota.NotifiedEventType,
		).
		Builder()
}
