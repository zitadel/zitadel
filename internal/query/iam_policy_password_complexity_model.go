package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
)

type IAMPasswordComplexityPolicyReadModel struct {
	PasswordComplexityPolicyReadModel
}

func (rm *IAMPasswordComplexityPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.PasswordComplexityPolicyAddedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *iam.PasswordComplexityPolicyChangedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		case *policy.PasswordComplexityPolicyAddedEvent,
			*policy.PasswordComplexityPolicyChangedEvent:

			rm.PasswordComplexityPolicyReadModel.AppendEvents(e)
		}
	}
}
