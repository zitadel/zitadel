package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type quotaNotificationsWriteModel struct {
	eventstore.WriteModel
	periodStart              time.Time
	latestNotifiedThresholds map[string]uint16
}

func newQuotaNotificationsWriteModel(aggregateId, instanceId, resourceOwner string, periodStart time.Time) *quotaNotificationsWriteModel {
	return &quotaNotificationsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   aggregateId,
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
		periodStart:              periodStart,
		latestNotifiedThresholds: make(map[string]uint16),
	}
}

func (wm *quotaNotificationsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		InstanceID(wm.InstanceID).
		AggregateTypes(quota.AggregateType).
		AggregateIDs(wm.AggregateID).
		CreationDateAfter(wm.periodStart).
		EventTypes(quota.NotifiedEventType).Builder()
}

func (wm *quotaNotificationsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		e := event.(*quota.NotifiedEvent)
		wm.latestNotifiedThresholds[e.ID] = e.Threshold
	}
	return wm.WriteModel.Reduce()
}
