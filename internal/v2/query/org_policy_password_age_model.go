package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/org"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type OrgPasswordAgePolicyReadModel struct {
	PasswordAgePolicyReadModel
}

func (rm *OrgPasswordAgePolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.PasswordAgePolicyAddedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(&e.PasswordAgePolicyAddedEvent)
		case *org.PasswordAgePolicyChangedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		case *policy.PasswordAgePolicyAddedEvent, *policy.PasswordAgePolicyChangedEvent:
			rm.PasswordAgePolicyReadModel.AppendEvents(e)
		}
	}
}
