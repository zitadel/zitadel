package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/settings"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMDebugNotificationLogWriteModel struct {
	DebugNotificationWriteModel
}

func NewIAMDebugNotificationLogWriteModel() *IAMDebugNotificationLogWriteModel {
	return &IAMDebugNotificationLogWriteModel{
		DebugNotificationWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
		},
	}
}

func (wm *IAMDebugNotificationLogWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.DebugNotificationProviderLogAddedEvent:
			wm.DebugNotificationWriteModel.AppendEvents(&e.DebugNotificationProviderAddedEvent)
		case *iam.DebugNotificationProviderLogChangedEvent:
			wm.DebugNotificationWriteModel.AppendEvents(&e.DebugNotificationProviderChangedEvent)
		case *iam.DebugNotificationProviderLogEnabledEvent:
			wm.DebugNotificationWriteModel.AppendEvents(&e.DebugNotificationProviderEnabledEvent)
		case *iam.DebugNotificationProviderLogDisabledEvent:
			wm.DebugNotificationWriteModel.AppendEvents(&e.DebugNotificationProviderDisabledEvent)
		case *iam.DebugNotificationProviderLogRemovedEvent:
			wm.DebugNotificationWriteModel.AppendEvents(&e.DebugNotificationProviderRemovedEvent)
		}
	}
}

func (wm *IAMDebugNotificationLogWriteModel) IsValid() bool {
	return wm.AggregateID != ""
}

func (wm *IAMDebugNotificationLogWriteModel) Reduce() error {
	return wm.DebugNotificationWriteModel.Reduce()
}

func (wm *IAMDebugNotificationLogWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.DebugNotificationWriteModel.AggregateID).
		EventTypes(
			iam.DebugNotificationProviderLogAddedEventType,
			iam.DebugNotificationProviderLogChangedEventType,
			iam.DebugNotificationProviderLogEnabledEventType,
			iam.DebugNotificationProviderLogDisabledEventType,
			iam.DebugNotificationProviderLogRemovedEventType).
		Builder()
}

func (wm *IAMDebugNotificationLogWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	compact bool) (*iam.DebugNotificationProviderLogChangedEvent, bool) {

	changes := make([]settings.DebugNotificationProviderChanges, 0)
	if wm.Compact != compact {
		changes = append(changes, settings.ChangeCompact(compact))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := iam.NewDebugNotificationProviderLogChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
