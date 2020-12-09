package password_age

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_age"
)

type PasswordAgePolicyReadModel struct {
	password_age.PasswordAgePolicyReadModel
}

func (rm *PasswordAgePolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *PasswordAgePolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordAgePolicyAddedEvent)
		case *PasswordAgePolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.PasswordAgePolicyChangedEvent)
		case *password_age.PasswordAgePolicyAddedEvent,
			*password_age.PasswordAgePolicyChangedEvent:

			rm.ReadModel.AppendEvents(e)
		}
	}
}
