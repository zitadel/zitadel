package idp

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

type ConfigsReadModel struct {
	eventstore.ReadModel

	Configs []*ConfigReadModel
}

func (rm *ConfigsReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		case *oidc.ConfigAddedEvent:
			rm.OIDCConfig = &oidc.ConfigReadModel{}
			rm.OIDCConfig.AppendEvents(event)
		case *oidc.ConfigChangedEvent:
			rm.OIDCConfig.AppendEvents(event)
		}
	}
}

func (rm *ConfigsReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *ConfigAddedEvent:
			rm.ConfigID = e.ConfigID
			rm.Name = e.Name
			rm.StylingType = e.StylingType
			rm.State = ConfigStateActive
		case *ConfigChangedEvent:
			if e.Name != "" {
				rm.Name = e.Name
			}
			if e.StylingType.Valid() {
				rm.StylingType = e.StylingType
			}
		case *ConfigDeactivatedEvent:
			rm.State = ConfigStateInactive
		case *ConfigReactivatedEvent:
			rm.State = ConfigStateActive
		case *ConfigRemovedEvent:
			rm.State = ConfigStateRemoved
		case *oidc.ConfigAddedEvent:
			rm.Type = ConfigTypeOIDC
		}
	}
	if err := rm.OIDCConfig.Reduce(); err != nil {
		return err
	}
	return rm.ReadModel.Reduce()
}
