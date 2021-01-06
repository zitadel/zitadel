package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type PasswordLockoutPolicyWriteModel struct {
	eventstore.WriteModel

	MaxAttempts         uint64
	ShowLockOutFailures bool
	IsActive            bool
}

func (wm *PasswordLockoutPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.PasswordLockoutPolicyAddedEvent:
			wm.MaxAttempts = e.MaxAttempts
			wm.ShowLockOutFailures = e.ShowLockOutFailures
			wm.IsActive = true
		case *policy.PasswordLockoutPolicyChangedEvent:
			if e.MaxAttempts != nil {
				wm.MaxAttempts = *e.MaxAttempts
			}
			if e.ShowLockOutFailures != nil {
				wm.ShowLockOutFailures = *e.ShowLockOutFailures
			}
		case *policy.PasswordLockoutPolicyRemovedEvent:
			wm.IsActive = false
		}
	}
	return wm.WriteModel.Reduce()
}
