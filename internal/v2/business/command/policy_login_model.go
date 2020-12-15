package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type LoginPolicyWriteModel struct {
	eventstore.WriteModel

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
	ForceMFA              bool
	PasswordlessType      domain.PasswordlessType
	IsActive              bool
}

func (wm *LoginPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.LoginPolicyAddedEvent:
			wm.AllowRegister = e.AllowRegister
			wm.AllowUserNamePassword = e.AllowUserNamePassword
			wm.AllowExternalIDP = e.AllowExternalIDP
			wm.ForceMFA = e.ForceMFA
			wm.PasswordlessType = e.PasswordlessType
			wm.IsActive = true
		case *policy.LoginPolicyChangedEvent:
			wm.AllowRegister = e.AllowRegister
			wm.AllowUserNamePassword = e.AllowUserNamePassword
			wm.AllowExternalIDP = e.AllowExternalIDP
			wm.ForceMFA = e.ForceMFA
			wm.PasswordlessType = e.PasswordlessType
		case *policy.LoginPolicyRemovedEvent:
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}
