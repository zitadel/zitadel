package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type PasswordComplexityPolicyReadModel struct {
	eventstore.ReadModel

	MinLength    uint64
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
}

func (rm *PasswordComplexityPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *policy.PasswordComplexityPolicyAddedEvent:
			rm.MinLength = e.MinLength
			rm.HasLowercase = e.HasLowercase
			rm.HasUpperCase = e.HasUppercase
			rm.HasNumber = e.HasNumber
			rm.HasSymbol = e.HasSymbol
		case *policy.PasswordComplexityPolicyChangedEvent:
			if e.MinLength != nil {
				rm.MinLength = *e.MinLength
			}
			if e.HasLowercase != nil {
				rm.HasLowercase = *e.HasLowercase
			}
			if e.HasUppercase != nil {
				rm.HasUpperCase = *e.HasUppercase
			}
			if e.HasNumber != nil {
				rm.HasNumber = *e.HasNumber
			}
			if e.HasSymbol != nil {
				rm.HasSymbol = *e.HasSymbol
			}
		}
	}
	return rm.ReadModel.Reduce()
}
