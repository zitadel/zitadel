package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/repository/instance"

	"github.com/zitadel/zitadel/internal/api/authz"

	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/eventstore"
)

type quotaWriteModel struct {
	eventstore.WriteModel
	unit   quota.Unit
	exists bool
}

func newQuotaWriteModel(ctx context.Context, unit quota.Unit) *quotaWriteModel {
	return &quotaWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   authz.GetInstance(ctx).InstanceID(),
			ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		},
		unit: unit,
	}
}

func (wm *quotaWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.QuotaAddedEvent:
			// TODO: Filter by quota aggregate id, then remove this check
			if e.Unit == wm.unit {
				wm.exists = true
			}
		case *instance.QuotaRemovedEvent:
			// TODO: Filter by quota aggregate id, then remove this check
			if e.Unit == wm.unit {
				wm.exists = false
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *quotaWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.QuotaAddedEventType,
			instance.QuotaRemovedEventType,
		).
		Builder()
}
