package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type OrgPasswordComplexityPolicyReadModel struct {
	PasswordComplexityPolicyReadModel
}

func (rm *OrgPasswordComplexityPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.PasswordComplexityPolicyAddedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *org.PasswordComplexityPolicyChangedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		case *policy.PasswordComplexityPolicyAddedEvent, *policy.PasswordComplexityPolicyChangedEvent:
			rm.PasswordComplexityPolicyReadModel.AppendEvents(e)
		}
	}
}
