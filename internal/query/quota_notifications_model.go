package query

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type quotaNotificationsReadModel struct {
	eventstore.ReadModel
	periodStart              time.Time
	latestNotifiedThresholds map[string]uint16
}

func newQuotaNotificationsReadModel(aggregateId, instanceId, resourceOwner string, periodStart time.Time) *quotaNotificationsReadModel {
	return &quotaNotificationsReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID:   aggregateId,
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
		periodStart:              periodStart,
		latestNotifiedThresholds: make(map[string]uint16),
	}
}

func (rm *quotaNotificationsReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(rm.ResourceOwner).
		AllowTimeTravel().
		AddQuery().
		InstanceID(rm.InstanceID).
		AggregateTypes(quota.AggregateType).
		AggregateIDs(rm.AggregateID).
		CreationDateAfter(rm.periodStart).
		EventTypes(quota.NotifiedEventType).Builder()
}

func (rm *quotaNotificationsReadModel) Reduce() error {
	for _, event := range rm.Events {
		e := event.(*quota.NotifiedEvent)
		rm.latestNotifiedThresholds[e.ID] = e.Threshold
	}
	return rm.ReadModel.Reduce()
}
