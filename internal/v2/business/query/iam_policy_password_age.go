package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IAMPasswordAgePolicyReadModel struct {
	PasswordAgePolicyReadModel
}

func (rm *IAMPasswordAgePolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordAgePolicyAddedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(&e.PassowordAgePolicyAddedEvent)
		case *iam.PasswordAgePolicyChangedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		case *policy.PassowordAgePolicyAddedEvent,
			*policy.PasswordAgePolicyChangedEvent:

			rm.PasswordAgePolicyReadModel.AppendEvents(e)
		}
	}
}
