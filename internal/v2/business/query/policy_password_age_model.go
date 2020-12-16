package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
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
		case *policy.PassowordAgePolicyAddedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		case *policy.PasswordAgePolicyChangedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		}
	}
	return rm.ReadModel.Reduce()
}
