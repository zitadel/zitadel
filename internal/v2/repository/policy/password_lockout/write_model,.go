package password_lockout

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type PasswordLockoutPolicyWriteModel struct {
	eventstore.WriteModel

	MaxAttempts         uint64
	ShowLockOutFailures bool
}

func (wm *PasswordLockoutPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *PasswordLockoutPolicyAddedEvent:
			wm.MaxAttempts = e.MaxAttempts
			wm.ShowLockOutFailures = e.ShowLockOutFailures
		case *PasswordLockoutPolicyChangedEvent:
			wm.MaxAttempts = e.MaxAttempts
			wm.ShowLockOutFailures = e.ShowLockOutFailures
		}
	}
	return wm.WriteModel.Reduce()
}
