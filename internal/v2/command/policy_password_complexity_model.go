package command

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type PasswordComplexityPolicyWriteModel struct {
	eventstore.WriteModel

	MinLength    uint64
	HasLowercase bool
	HasUppercase bool
	HasNumber    bool
	HasSymbol    bool
	State        domain.PolicyState
}

func (wm *PasswordComplexityPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.PasswordComplexityPolicyAddedEvent:
			wm.MinLength = e.MinLength
			wm.HasLowercase = e.HasLowercase
			wm.HasUppercase = e.HasUppercase
			wm.HasNumber = e.HasNumber
			wm.HasSymbol = e.HasSymbol
			wm.State = domain.PolicyStateActive
		case *policy.PasswordComplexityPolicyChangedEvent:
			if e.MinLength != nil {
				wm.MinLength = *e.MinLength
			}
			if e.HasLowercase != nil {
				wm.HasLowercase = *e.HasLowercase
			}
			if e.HasUppercase != nil {
				wm.HasUppercase = *e.HasUppercase
			}
			if e.HasNumber != nil {
				wm.HasNumber = *e.HasNumber
			}
			if e.HasSymbol != nil {
				wm.HasSymbol = *e.HasSymbol
			}
		case *policy.PasswordComplexityPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
