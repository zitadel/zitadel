package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
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
