package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type LockoutPolicyWriteModel struct {
	eventstore.WriteModel

	MaxPasswordAttempts uint64
	MaxOTPAttempts      uint64
	ShowLockOutFailures bool
	State               domain.PolicyState
}

func (wm *LockoutPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.LockoutPolicyAddedEvent:
			wm.MaxPasswordAttempts = e.MaxPasswordAttempts
			wm.MaxOTPAttempts = e.MaxOTPAttempts
			wm.ShowLockOutFailures = e.ShowLockOutFailures
			wm.State = domain.PolicyStateActive
		case *policy.LockoutPolicyChangedEvent:
			if e.MaxPasswordAttempts != nil {
				wm.MaxPasswordAttempts = *e.MaxPasswordAttempts
			}
			if e.MaxOTPAttempts != nil {
				wm.MaxOTPAttempts = *e.MaxOTPAttempts
			}
			if e.ShowLockOutFailures != nil {
				wm.ShowLockOutFailures = *e.ShowLockOutFailures
			}
		case *policy.LockoutPolicyRemovedEvent:
			wm.State = domain.PolicyStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
