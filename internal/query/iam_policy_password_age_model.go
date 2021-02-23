package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMPasswordAgePolicyReadModel struct {
	PasswordAgePolicyReadModel
}

func (rm *IAMPasswordAgePolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordAgePolicyAddedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(&e.PasswordAgePolicyAddedEvent)
		case *iam.PasswordAgePolicyChangedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		case *policy.PasswordAgePolicyAddedEvent,
			*policy.PasswordAgePolicyChangedEvent:

			rm.PasswordAgePolicyReadModel.AppendEvents(e)
		}
	}
}
