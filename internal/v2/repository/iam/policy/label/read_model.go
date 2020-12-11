package label

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/label"
)

type ReadModel struct{ label.ReadModel }

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *ChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *label.AddedEvent, *label.ChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}
