package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type PasswordLockoutPolicyWriteModel struct {
	eventstore.WriteModel

	MaxAttempts         uint64
	ShowLockOutFailures bool
	State               domain.PolicyState
}

func (wm *PasswordLockoutPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.PasswordLockoutPolicyAddedEvent:
			wm.MaxAttempts = e.MaxAttempts
			wm.ShowLockOutFailures = e.ShowLockOutFailures
			wm.State = domain.PolicyStateActive
		case *policy.PasswordLockoutPolicyChangedEvent:
			if e.MaxAttempts != nil {
				wm.MaxAttempts = *e.MaxAttempts
			}
			if e.ShowLockOutFailures != nil {
				wm.ShowLockOutFailures = *e.ShowLockOutFailures
			}
		case *policy.PasswordLockoutPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
