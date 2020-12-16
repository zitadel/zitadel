package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type PasswordAgePolicyWriteModel struct {
	eventstore.WriteModel

	ExpireWarnDays uint64
	MaxAgeDays     uint64
	IsActive       bool
}

func (wm *PasswordAgePolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.PassowordAgePolicyAddedEvent:
			wm.ExpireWarnDays = e.ExpireWarnDays
			wm.MaxAgeDays = e.MaxAgeDays
			wm.IsActive = true
		case *policy.PasswordAgePolicyChangedEvent:
			wm.ExpireWarnDays = e.ExpireWarnDays
			wm.MaxAgeDays = e.MaxAgeDays
		case *policy.PasswordAgePolicyRemovedEvent:
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}
