package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type PasswordAgePolicyReadModel struct {
	eventstore.ReadModel

	ExpireWarnDays uint64
	MaxAgeDays     uint64
}

func (rm *PasswordAgePolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *policy.PasswordAgePolicyAddedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		case *policy.PasswordAgePolicyChangedEvent:
			if e.ExpireWarnDays != nil {
				rm.ExpireWarnDays = *e.ExpireWarnDays
			}
			if e.MaxAgeDays != nil {
				rm.MaxAgeDays = *e.MaxAgeDays
			}
		}
	}
	return rm.ReadModel.Reduce()
}
