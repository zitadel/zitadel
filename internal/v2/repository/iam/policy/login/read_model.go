package login

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
)

type LoginPolicyReadModel struct{ login.LoginPolicyReadModel }

func (rm *LoginPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LoginPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.LoginPolicyAddedEvent)
		case *LoginPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.LoginPolicyChangedEvent)
		case *login.LoginPolicyAddedEvent, *login.LoginPolicyChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}
