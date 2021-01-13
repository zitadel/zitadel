package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
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
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, org.AggregateType).
		AggregateIDs(wm.LoginPolicyWriteModel.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *OrgLoginPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType domain.PasswordlessType,
) (*org.LoginPolicyChangedEvent, bool) {

	hasChanged := false
	changedEvent := org.NewLoginPolicyChangedEvent(ctx)
	if wm.AllowUserNamePassword == allowUsernamePassword {
		hasChanged = true
		changedEvent.AllowUserNamePassword = &allowUsernamePassword
	}
	if wm.AllowRegister == allowRegister {
		hasChanged = true
		changedEvent.AllowRegister = &allowRegister
	}
	if wm.AllowExternalIDP == allowExternalIDP {
		hasChanged = true
		changedEvent.AllowExternalIDP = &allowExternalIDP
	}
	if wm.ForceMFA != forceMFA {
		hasChanged = true
		changedEvent.ForceMFA = &forceMFA
	}
	if passwordlessType.Valid() && wm.PasswordlessType != passwordlessType {
		hasChanged = true
		changedEvent.PasswordlessType = &passwordlessType
	}
	return changedEvent, hasChanged
}
