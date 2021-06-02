package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type OrgLoginPolicyWriteModel struct {
	LoginPolicyWriteModel
}

func NewOrgLoginPolicyWriteModel(orgID string) *OrgLoginPolicyWriteModel {
	return &OrgLoginPolicyWriteModel{
		LoginPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
		},
	}
}

func (wm *OrgLoginPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LoginPolicyAddedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *org.LoginPolicyChangedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyChangedEvent)
		case *org.LoginPolicyRemovedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyRemovedEvent)
		}
	}
}

func (wm *OrgLoginPolicyWriteModel) IsValid() bool {
	return wm.AggregateID != ""
}

func (wm *OrgLoginPolicyWriteModel) Reduce() error {
	return wm.LoginPolicyWriteModel.Reduce()
}

func (wm *OrgLoginPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.LoginPolicyWriteModel.AggregateID).
		EventTypes(
			org.LoginPolicyAddedEventType,
			org.LoginPolicyChangedEventType,
			org.LoginPolicyRemovedEventType).
		Builder()
}

func (wm *OrgLoginPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType domain.PasswordlessType,
) (*org.LoginPolicyChangedEvent, bool) {

	changes := make([]policy.LoginPolicyChanges, 0)
	if wm.AllowUserNamePassword != allowUsernamePassword {
		changes = append(changes, policy.ChangeAllowUserNamePassword(allowUsernamePassword))
	}
	if wm.AllowRegister != allowRegister {
		changes = append(changes, policy.ChangeAllowRegister(allowRegister))
	}
	if wm.AllowExternalIDP != allowExternalIDP {
		changes = append(changes, policy.ChangeAllowExternalIDP(allowExternalIDP))
	}
	if wm.ForceMFA != forceMFA {
		changes = append(changes, policy.ChangeForceMFA(forceMFA))
	}
	if passwordlessType.Valid() && wm.PasswordlessType != passwordlessType {
		changes = append(changes, policy.ChangePasswordlessType(passwordlessType))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := org.NewLoginPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
