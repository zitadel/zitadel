package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type OrgNotificationPolicyWriteModel struct {
	NotificationPolicyWriteModel
}

func NewOrgNotificationPolicyWriteModel(orgID string) *OrgNotificationPolicyWriteModel {
	return &OrgNotificationPolicyWriteModel{
		NotificationPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgNotificationPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.NotificationPolicyAddedEvent:
			wm.NotificationPolicyWriteModel.AppendEvents(&e.NotificationPolicyAddedEvent)
		case *org.NotificationPolicyChangedEvent:
			wm.NotificationPolicyWriteModel.AppendEvents(&e.NotificationPolicyChangedEvent)
		case *org.NotificationPolicyRemovedEvent:
			wm.NotificationPolicyWriteModel.AppendEvents(&e.NotificationPolicyRemovedEvent)
		}
	}
}

func (wm *OrgNotificationPolicyWriteModel) Reduce() error {
	return wm.NotificationPolicyWriteModel.Reduce()
}

func (wm *OrgNotificationPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.NotificationPolicyWriteModel.AggregateID).
		AggregateTypes(org.AggregateType).
		EventTypes(org.NotificationPolicyAddedEventType,
			org.NotificationPolicyChangedEventType,
			org.NotificationPolicyRemovedEventType).
		Builder()
}

func (wm *OrgNotificationPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	passwordChange bool,
) (*org.NotificationPolicyChangedEvent, bool) {

	changes := make([]policy.NotificationPolicyChanges, 0)
	if wm.PasswordChange != passwordChange {
		changes = append(changes, policy.ChangePasswordChange(passwordChange))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewNotificationPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
