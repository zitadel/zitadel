package password_age

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type PasswordAgePolicyWriteModel struct {
	eventstore.WriteModel

	ExpireWarnDays uint64
	MaxAgeDays     uint64
}

func (wm *PasswordAgePolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *PasswordAgePolicyAddedEvent:
			wm.ExpireWarnDays = e.ExpireWarnDays
			wm.MaxAgeDays = e.MaxAgeDays
		case *PasswordAgePolicyChangedEvent:
			wm.ExpireWarnDays = e.ExpireWarnDays
			wm.MaxAgeDays = e.MaxAgeDays
		}
	}
	return wm.WriteModel.Reduce()
}
