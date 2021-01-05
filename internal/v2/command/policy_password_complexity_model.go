package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type PasswordComplexityPolicyWriteModel struct {
	eventstore.WriteModel

	MinLength    uint64
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
	IsActive     bool
}

func (wm *PasswordComplexityPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.PasswordComplexityPolicyAddedEvent:
			wm.MinLength = e.MinLength
			wm.HasLowercase = e.HasLowercase
			wm.HasUpperCase = e.HasUpperCase
			wm.HasNumber = e.HasNumber
			wm.HasSymbol = e.HasSymbol
			wm.IsActive = true
		case *policy.PasswordComplexityPolicyChangedEvent:
			if e.MinLength != nil {
				wm.MinLength = *e.MinLength
			}
			if e.HasLowercase != nil {
				wm.HasLowercase = *e.HasLowercase
			}
			if e.HasUpperCase != nil {
				wm.HasUpperCase = *e.HasUpperCase
			}
			if e.HasNumber != nil {
				wm.HasNumber = *e.HasNumber
			}
			if e.HasSymbol != nil {
				wm.HasSymbol = *e.HasSymbol
			}
		case *policy.PasswordComplexityPolicyRemovedEvent:
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}
