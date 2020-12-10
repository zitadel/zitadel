package password_lockout

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_lockout"
)

type ReadModel struct {
	password_lockout.ReadModel
}

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *ChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *password_lockout.AddedEvent, *password_lockout.ChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}
