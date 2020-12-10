package password_age

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type WriteModel struct {
	eventstore.WriteModel

	ExpireWarnDays uint64
	MaxAgeDays     uint64
}

func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.ExpireWarnDays = e.ExpireWarnDays
			wm.MaxAgeDays = e.MaxAgeDays
		case *ChangedEvent:
			wm.ExpireWarnDays = e.ExpireWarnDays
			wm.MaxAgeDays = e.MaxAgeDays
		}
	}
	return wm.WriteModel.Reduce()
}
