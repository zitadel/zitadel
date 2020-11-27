package idp

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

type ConfigWriteModel struct {
	eventstore.WriteModel

	State ConfigState

	ConfigID    string
	Name        string
	StylingType StylingType

	OIDCConfig *oidc.ConfigWriteModel
}

func (rm *ConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	rm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		case *oidc.ConfigAddedEvent:
			rm.OIDCConfig = new(oidc.ConfigWriteModel)
			rm.OIDCConfig.AppendEvents(event)
		case *oidc.ConfigChangedEvent:
			rm.OIDCConfig.AppendEvents(event)
		}
	}
}

func (rm *ConfigWriteModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *ConfigAddedEvent:
			rm.reduceConfigAddedEvent(e)
		case *ConfigChangedEvent:
			rm.reduceConfigChangedEvent(e)
		case *ConfigDeactivatedEvent:
			rm.reduceConfigStateChanged(e.ConfigID, ConfigStateInactive)
		case *ConfigReactivatedEvent:
			rm.reduceConfigStateChanged(e.ConfigID, ConfigStateActive)
		case *ConfigRemovedEvent:
			rm.reduceConfigStateChanged(e.ConfigID, ConfigStateRemoved)
		}
	}
	if rm.OIDCConfig != nil {
		if err := rm.OIDCConfig.Reduce(); err != nil {
			return err
		}
	}
	return rm.WriteModel.Reduce()
}

func (rm *ConfigWriteModel) reduceConfigAddedEvent(e *ConfigAddedEvent) {
	rm.ConfigID = e.ConfigID
	rm.Name = e.Name
	rm.StylingType = e.StylingType
	rm.State = ConfigStateActive
}

func (rm *ConfigWriteModel) reduceConfigChangedEvent(e *ConfigChangedEvent) {
	if e.Name != "" {
		rm.Name = e.Name
	}
	if e.StylingType.Valid() {
		rm.StylingType = e.StylingType
	}
}

func (rm *ConfigWriteModel) reduceConfigStateChanged(configID string, state ConfigState) {
	rm.State = state
}
