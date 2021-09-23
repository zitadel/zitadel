package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMLockoutPolicyReadModel struct {
	LockoutPolicyReadModel
}

func (rm *IAMLockoutPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LockoutPolicyAddedEvent:
			rm.LockoutPolicyReadModel.AppendEvents(&e.LockoutPolicyAddedEvent)
		case *iam.LockoutPolicyChangedEvent:
			rm.LockoutPolicyReadModel.AppendEvents(&e.LockoutPolicyChangedEvent)
		case *policy.LockoutPolicyAddedEvent, *policy.LockoutPolicyChangedEvent:
			rm.LockoutPolicyReadModel.AppendEvents(e)
		}
	}
}
