package password_complexity

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_complexity"
)

type PasswordComplexityPolicyReadModel struct {
	password_complexity.PasswordComplexityPolicyReadModel
}

func (rm *PasswordComplexityPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordComplexityPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *PasswordComplexityPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		case *password_complexity.PasswordComplexityPolicyAddedEvent,
			*password_complexity.PasswordComplexityPolicyChangedEvent:

			rm.ReadModel.AppendEvents(e)
		}
	}
}
