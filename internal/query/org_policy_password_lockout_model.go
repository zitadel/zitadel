package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

type OrgPasswordLockoutPolicyReadModel struct {
	LockoutPolicyReadModel
}

func (rm *OrgPasswordLockoutPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.LockoutPolicyAddedEvent:
			rm.LockoutPolicyReadModel.AppendEvents(&e.LockoutPolicyAddedEvent)
		case *org.LockoutPolicyChangedEvent:
			rm.LockoutPolicyReadModel.AppendEvents(&e.LockoutPolicyChangedEvent)
		case *policy.LockoutPolicyAddedEvent, *policy.LockoutPolicyChangedEvent:
			rm.LockoutPolicyReadModel.AppendEvents(e)
		}
	}
}
