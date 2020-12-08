package idpprovider

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

type IDPProviderWriteModel struct {
	provider.WriteModel
}

func (wm *IDPProviderWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *IDPProviderAddedEvent:
			wm.WriteModel.AppendEvents(&e.AddedEvent)
		}
	}
}
