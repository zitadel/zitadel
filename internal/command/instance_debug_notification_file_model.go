package command

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/settings"

	"github.com/caos/zitadel/internal/repository/instance"
)

type InstanceDebugNotificationFileWriteModel struct {
	DebugNotificationWriteModel
}

func NewInstanceDebugNotificationFileWriteModel(ctx context.Context) *InstanceDebugNotificationFileWriteModel {
	return &InstanceDebugNotificationFileWriteModel{
		DebugNotificationWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstanceDebugNotificationFileWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.DebugNotificationProviderFileAddedEvent:
			wm.DebugNotificationWriteModel.AppendEvents(&e.DebugNotificationProviderAddedEvent)
		case *instance.DebugNotificationProviderFileChangedEvent:
			wm.DebugNotificationWriteModel.AppendEvents(&e.DebugNotificationProviderChangedEvent)
		case *instance.DebugNotificationProviderFileRemovedEvent:
			wm.DebugNotificationWriteModel.AppendEvents(&e.DebugNotificationProviderRemovedEvent)
		}
	}
}

func (wm *InstanceDebugNotificationFileWriteModel) IsValid() bool {
	return wm.AggregateID != ""
}

func (wm *InstanceDebugNotificationFileWriteModel) Reduce() error {
	return wm.DebugNotificationWriteModel.Reduce()
}

func (wm *InstanceDebugNotificationFileWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.DebugNotificationWriteModel.AggregateID).
		EventTypes(
			instance.DebugNotificationProviderFileAddedEventType,
			instance.DebugNotificationProviderFileChangedEventType,
			instance.DebugNotificationProviderFileRemovedEventType).
		Builder()
}

func (wm *InstanceDebugNotificationFileWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	compact bool) (*instance.DebugNotificationProviderFileChangedEvent, bool) {

	changes := make([]settings.DebugNotificationProviderChanges, 0)
	if wm.Compact != compact {
		changes = append(changes, settings.ChangeCompact(compact))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewDebugNotificationProviderFileChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
