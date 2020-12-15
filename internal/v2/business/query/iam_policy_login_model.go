package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

type IAMLoginPolicyReadModel struct{ LoginPolicyReadModel }

func (rm *IAMLoginPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.LoginPolicyAddedEvent:
			rm.LoginPolicyReadModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *iam.LoginPolicyChangedEvent:
			rm.LoginPolicyReadModel.AppendEvents(&e.LoginPolicyChangedEvent)
		case *policy.LoginPolicyAddedEvent, *policy.LoginPolicyChangedEvent:
			rm.LoginPolicyReadModel.AppendEvents(e)
		}
	}
}
