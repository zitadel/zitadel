package login

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type LoginPolicyWriteModel struct {
	eventstore.WriteModel

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
	ForceMFA              bool
	PasswordlessType      PasswordlessType
}

func (wm *LoginPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			wm.AllowRegister = e.AllowRegister
			wm.AllowUserNamePassword = e.AllowUserNamePassword
			wm.AllowExternalIDP = e.AllowExternalIDP
			wm.ForceMFA = e.ForceMFA
			wm.PasswordlessType = e.PasswordlessType
		case *LoginPolicyChangedEvent:
			wm.AllowRegister = e.AllowRegister
			wm.AllowUserNamePassword = e.AllowUserNamePassword
			wm.AllowExternalIDP = e.AllowExternalIDP
			wm.ForceMFA = e.ForceMFA
			wm.PasswordlessType = e.PasswordlessType
		}
	}
	return wm.WriteModel.Reduce()
}
