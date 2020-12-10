package login

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type ReadModel struct {
	eventstore.ReadModel

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
	ForceMFA              bool
	PasswordlessType      PasswordlessType
}

func (rm *ReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.AllowUserNamePassword = e.AllowUserNamePassword
			rm.AllowExternalIDP = e.AllowExternalIDP
			rm.AllowRegister = e.AllowRegister
			rm.ForceMFA = e.ForceMFA
			rm.PasswordlessType = e.PasswordlessType
		case *ChangedEvent:
			rm.AllowUserNamePassword = e.AllowUserNamePassword
			rm.AllowExternalIDP = e.AllowExternalIDP
			rm.AllowRegister = e.AllowRegister
			rm.ForceMFA = e.ForceMFA
			rm.PasswordlessType = e.PasswordlessType
		}
	}
	return rm.ReadModel.Reduce()
}
