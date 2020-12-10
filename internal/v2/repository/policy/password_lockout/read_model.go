package password_lockout

import "github.com/caos/zitadel/internal/eventstore/v2"

type ReadModel struct {
	eventstore.ReadModel

	MaxAttempts         uint64
	ShowLockOutFailures bool
}

func (rm *ReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.MaxAttempts = e.MaxAttempts
			rm.ShowLockOutFailures = e.ShowLockOutFailures
		case *ChangedEvent:
			rm.MaxAttempts = e.MaxAttempts
			rm.ShowLockOutFailures = e.ShowLockOutFailures
		}
	}
	return rm.ReadModel.Reduce()
}
