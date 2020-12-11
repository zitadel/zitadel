package login

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
)

type ReadModel struct{ login.ReadModel }

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *ChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *login.AddedEvent, *login.ChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}
