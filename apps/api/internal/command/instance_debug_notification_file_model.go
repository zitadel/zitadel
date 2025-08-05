package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/settings"
)

type InstanceDebugNotificationFileWriteModel struct {
	DebugNotificationWriteModel
}

func NewInstanceDebugNotificationFileWriteModel(ctx context.Context) *InstanceDebugNotificationFileWriteModel {
	instanceID := authz.GetInstance(ctx).InstanceID()
	return &InstanceDebugNotificationFileWriteModel{
		DebugNotificationWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
				InstanceID:    instanceID,
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
