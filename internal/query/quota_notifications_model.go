package query

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type quotaNotificationsReadModel struct {
	eventstore.ReadModel
	periodStart         time.Time
	latestDueThresholds map[string]uint16
}

func newQuotaNotificationsReadModel(aggregateId, instanceId, resourceOwner string, periodStart time.Time) *quotaNotificationsReadModel {
	return &quotaNotificationsReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID:   aggregateId,
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
		periodStart:         periodStart,
		latestDueThresholds: make(map[string]uint16),
	}
}

func (rm *quotaNotificationsReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		ResourceOwner(rm.ResourceOwner).
		AllowTimeTravel().
		AddQuery().
		InstanceID(rm.InstanceID).
		AggregateTypes(quota.AggregateType).
		AggregateIDs(rm.AggregateID).
		CreationDateAfter(rm.periodStart).
		EventTypes(quota.NotificationDueEventType).Builder()
}

func (rm *quotaNotificationsReadModel) Reduce() error {
	for _, event := range rm.Events {
		e := event.(*quota.NotificationDueEvent)
		rm.latestDueThresholds[e.ID] = e.Threshold
	}
	return rm.ReadModel.Reduce()
}
