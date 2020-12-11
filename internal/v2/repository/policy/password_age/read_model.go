package password_age

import "github.com/caos/zitadel/internal/eventstore/v2"

type ReadModel struct {
	eventstore.ReadModel

	ExpireWarnDays uint64
	MaxAgeDays     uint64
}

func (rm *ReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		case *ChangedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		}
	}
	return rm.ReadModel.Reduce()
}
