package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMLoginPolicyWriteModel struct {
	LoginPolicyWriteModel
}

func NewIAMLoginPolicyWriteModel(iamID string) *IAMLoginPolicyWriteModel {
	return &IAMLoginPolicyWriteModel{
		LoginPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
		},
	}
}

func (wm *IAMLoginPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LoginPolicyAddedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *iam.LoginPolicyChangedEvent:
			wm.LoginPolicyWriteModel.AppendEvents(&e.LoginPolicyChangedEvent)
		}
	}
}

func (wm *IAMLoginPolicyWriteModel) Reduce() error {
	return wm.LoginPolicyWriteModel.Reduce()
}

func (wm *IAMLoginPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.LoginPolicyWriteModel.AggregateID)
}

func (wm *IAMLoginPolicyWriteModel) NewChangedEvent(allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA bool, passwordlessType domain.PasswordlessType) (*iam.LoginPolicyChangedEvent, bool) {
	hasChanged := false
	changedEvent := &iam.LoginPolicyChangedEvent{}
	if wm.AllowUserNamePassword == allowUsernamePassword {
		hasChanged = true
		changedEvent.AllowUserNamePassword = allowUsernamePassword
	}
	if wm.AllowRegister == allowRegister {
		hasChanged = true
		changedEvent.AllowRegister = allowRegister
	}
	if wm.AllowExternalIDP == allowExternalIDP {
		hasChanged = true
		changedEvent.AllowExternalIDP = allowExternalIDP
	}
	if wm.ForceMFA != forceMFA {
		hasChanged = true
		changedEvent.ForceMFA = forceMFA
	}
	if wm.PasswordlessType != passwordlessType {
		hasChanged = true
		changedEvent.PasswordlessType = passwordlessType
	}
	return changedEvent, hasChanged
}
