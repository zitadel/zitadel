package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type LoginPolicyReadModel struct {
	eventstore.ReadModel

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
	ForceMFA              bool
	PasswordlessType      domain.PasswordlessType
	IsActive              bool
}

func (rm *LoginPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *policy.LoginPolicyAddedEvent:
			rm.AllowUserNamePassword = e.AllowUserNamePassword
			rm.AllowExternalIDP = e.AllowExternalIDP
			rm.AllowRegister = e.AllowRegister
			rm.ForceMFA = e.ForceMFA
			rm.PasswordlessType = e.PasswordlessType
			rm.IsActive = true
		case *policy.LoginPolicyChangedEvent:
			if e.AllowUserNamePassword != nil {
				rm.AllowUserNamePassword = *e.AllowUserNamePassword
			}
			if e.AllowExternalIDP != nil {
				rm.AllowExternalIDP = *e.AllowExternalIDP
			}
			if e.AllowRegister != nil {
				rm.AllowRegister = *e.AllowRegister
			}
			if e.ForceMFA != nil {
				rm.ForceMFA = *e.ForceMFA
			}
			if e.PasswordlessType != nil {
				rm.PasswordlessType = *e.PasswordlessType
			}
		case *policy.LoginPolicyRemovedEvent:
			rm.IsActive = false
		}
	}
	return rm.ReadModel.Reduce()
}
