package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IAMPasswordLockoutPolicyReadModel struct {
	PasswordLockoutPolicyReadModel
}

func (rm *IAMPasswordLockoutPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordLockoutPolicyAddedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(&e.PasswordLockoutPolicyAddedEvent)
		case *iam.PasswordLockoutPolicyChangedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(&e.PasswordLockoutPolicyChangedEvent)
		case *policy.PasswordLockoutPolicyAddedEvent, *policy.PasswordLockoutPolicyChangedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(e)
		}
	}
}
