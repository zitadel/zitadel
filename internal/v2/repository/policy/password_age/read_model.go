package password_age

import "github.com/caos/zitadel/internal/eventstore/v2"

type PasswordAgePolicyReadModel struct {
	eventstore.ReadModel

	ExpireWarnDays uint64
	MaxAgeDays     uint64
}

func (rm *PasswordAgePolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *PasswordAgePolicyAddedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		case *PasswordAgePolicyChangedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		}
	}
	return rm.ReadModel.Reduce()
}
