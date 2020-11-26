package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp"
)

type IDPConfigsReadModel struct {
	idp.ConfigsReadModel
}

func (rm *IDPConfigsReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *IDPConfigAddedEvent:
			rm.ConfigsReadModel.AppendEvents(&e.ConfigAddedEvent)
		case *IDPConfigChangedEvent:
			rm.ConfigsReadModel.AppendEvents(&e.ConfigChangedEvent)
		case *IDPConfigDeactivatedEvent:
			rm.ConfigsReadModel.AppendEvents(&e.ConfigDeactivatedEvent)
		case *IDPConfigReactivatedEvent:
			rm.ConfigsReadModel.AppendEvents(&e.ConfigReactivatedEvent)
		case *IDPConfigRemovedEvent:
			rm.ConfigsReadModel.AppendEvents(&e.ConfigRemovedEvent)
		case *IDPOIDCConfigAddedEvent:
			rm.ConfigsReadModel.AppendEvents(&e.ConfigAddedEvent)
		case *IDPOIDCConfigChangedEvent:
			rm.ConfigsReadModel.AppendEvents(&e.ConfigChangedEvent)
		}
	}
}
