package idpprovider

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

type WriteModel struct {
	provider.WriteModel
}

func (wm *WriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *AddedEvent:
			wm.WriteModel.AppendEvents(&e.AddedEvent)
		}
	}
}
