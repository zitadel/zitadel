package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMLoginPolicyWriteModel struct {
	LoginPolicyWriteModel
}

func NewWriteModel(iamID string) *IAMLoginPolicyWriteModel {
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

func (wm *IAMLoginPolicyWriteModel) HasChanged(allowUsernamePassword, allowRegister, allowExternalIDP, forceMFA bool, passwordlessType domain.PasswordlessType) bool {
	if wm.AllowUserNamePassword != allowUsernamePassword {
		return true
	}
	if wm.AllowRegister != allowRegister {
		return true
	}
	if wm.AllowExternalIDP != allowExternalIDP {
		return true
	}
	if wm.ForceMFA != forceMFA {
		return true
	}
	if wm.PasswordlessType != passwordlessType {
		return true
	}
	return false
}
