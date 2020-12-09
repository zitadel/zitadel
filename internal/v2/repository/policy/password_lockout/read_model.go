package password_lockout

import "github.com/caos/zitadel/internal/eventstore/v2"

type PasswordLockoutPolicyReadModel struct {
	eventstore.ReadModel

	MaxAttempts         uint64
	ShowLockOutFailures bool
}

func (rm *PasswordLockoutPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *PasswordLockoutPolicyAddedEvent:
			rm.MaxAttempts = e.MaxAttempts
			rm.ShowLockOutFailures = e.ShowLockOutFailures
		case *PasswordLockoutPolicyChangedEvent:
			rm.MaxAttempts = e.MaxAttempts
			rm.ShowLockOutFailures = e.ShowLockOutFailures
		}
	}
	return rm.ReadModel.Reduce()
}
