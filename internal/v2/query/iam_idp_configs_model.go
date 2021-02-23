package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMIDPConfigsReadModel struct {
	IDPConfigsReadModel
}

func (rm *IAMIDPConfigsReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.IDPConfigAddedEvent:
			rm.IDPConfigsReadModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *iam.IDPConfigChangedEvent:
			rm.IDPConfigsReadModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *iam.IDPConfigDeactivatedEvent:
			rm.IDPConfigsReadModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *iam.IDPConfigReactivatedEvent:
			rm.IDPConfigsReadModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *iam.IDPConfigRemovedEvent:
			rm.IDPConfigsReadModel.AppendEvents(&e.IDPConfigRemovedEvent)
		case *iam.IDPOIDCConfigAddedEvent:
			rm.IDPConfigsReadModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *iam.IDPOIDCConfigChangedEvent:
			rm.IDPConfigsReadModel.AppendEvents(&e.OIDCConfigChangedEvent)
		}
	}
}
