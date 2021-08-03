package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/policy"
)

type PasswordLockoutPolicyReadModel struct {
	eventstore.ReadModel

	MaxAttempts         uint64
	ShowLockOutFailures bool
}

func (rm *PasswordLockoutPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *policy.LockoutPolicyAddedEvent:
			rm.MaxAttempts = e.MaxPasswordAttempts
			rm.ShowLockOutFailures = e.ShowLockOutFailures
		case *policy.LockoutPolicyChangedEvent:
			if e.MaxPasswordAttempts != nil {
				rm.MaxAttempts = *e.MaxPasswordAttempts
			}
			if e.ShowLockOutFailures != nil {
				rm.ShowLockOutFailures = *e.ShowLockOutFailures
			}
		}
	}
	return rm.ReadModel.Reduce()
}
