package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type OrgPasswordLockoutPolicyReadModel struct {
	PasswordLockoutPolicyReadModel
}

func (rm *OrgPasswordLockoutPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.PasswordLockoutPolicyAddedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(&e.PasswordLockoutPolicyAddedEvent)
		case *org.PasswordLockoutPolicyChangedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(&e.PasswordLockoutPolicyChangedEvent)
		case *policy.PasswordLockoutPolicyAddedEvent, *policy.PasswordLockoutPolicyChangedEvent:
			rm.PasswordLockoutPolicyReadModel.AppendEvents(e)
		}
	}
}
