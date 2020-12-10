package password_lockout

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type WriteModel struct {
	eventstore.WriteModel

	MaxAttempts         uint64
	ShowLockOutFailures bool
}

func (wm *WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.MaxAttempts = e.MaxAttempts
			wm.ShowLockOutFailures = e.ShowLockOutFailures
		case *ChangedEvent:
			wm.MaxAttempts = e.MaxAttempts
			wm.ShowLockOutFailures = e.ShowLockOutFailures
		}
	}
	return wm.WriteModel.Reduce()
}
