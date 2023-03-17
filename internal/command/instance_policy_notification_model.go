package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type InstanceNotificationPolicyWriteModel struct {
	NotificationPolicyWriteModel
}

func NewInstanceNotificationPolicyWriteModel(ctx context.Context) *InstanceNotificationPolicyWriteModel {
	return &InstanceNotificationPolicyWriteModel{
		NotificationPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
		},
	}
}

func (wm *InstanceNotificationPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.NotificationPolicyAddedEvent:
			wm.NotificationPolicyWriteModel.AppendEvents(&e.NotificationPolicyAddedEvent)
		case *instance.NotificationPolicyChangedEvent:
			wm.NotificationPolicyWriteModel.AppendEvents(&e.NotificationPolicyChangedEvent)
		}
	}
}

func (wm *InstanceNotificationPolicyWriteModel) Reduce() error {
	return wm.NotificationPolicyWriteModel.Reduce()
}

func (wm *InstanceNotificationPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.NotificationPolicyWriteModel.AggregateID).
		EventTypes(
			instance.NotificationPolicyAddedEventType,
			instance.NotificationPolicyChangedEventType).
		Builder()
}

func (wm *InstanceNotificationPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	passwordChange bool,
) (*instance.NotificationPolicyChangedEvent, bool) {

	changes := make([]policy.NotificationPolicyChanges, 0)
	if wm.PasswordChange != passwordChange {
		changes = append(changes, policy.ChangePasswordChange(passwordChange))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewNotificationPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
