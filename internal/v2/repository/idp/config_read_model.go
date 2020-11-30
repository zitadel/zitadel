package idp

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

type ConfigReadModel struct {
	eventstore.ReadModel

	State        ConfigState
	ConfigID     string
	Name         string
	StylingType  StylingType
	ProviderType provider.Type

	OIDCConfig *oidc.ConfigReadModel
}

func NewConfigReadModel(configID string) *ConfigReadModel {
	return &ConfigReadModel{
		ConfigID: configID,
	}
}

func (rm *ConfigReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *ConfigAddedEvent:
			rm.ReadModel.AppendEvents(e)
		case *ConfigChangedEvent:
			rm.ReadModel.AppendEvents(e)
		case *ConfigDeactivatedEvent:
			rm.ReadModel.AppendEvents(e)
		case *ConfigReactivatedEvent:
			rm.ReadModel.AppendEvents(e)
		case *ConfigRemovedEvent:
			rm.ReadModel.AppendEvents(e)
		case *oidc.ConfigAddedEvent:
			rm.OIDCConfig = &oidc.ConfigReadModel{}
			rm.ReadModel.AppendEvents(e)
			rm.OIDCConfig.AppendEvents(event)
		case *oidc.ConfigChangedEvent:
			rm.ReadModel.AppendEvents(e)
			rm.OIDCConfig.AppendEvents(event)
		}
	}
}

func (rm *ConfigReadModel) Reduce() error {
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
	return rm.ReadModel.Reduce()
}

func (rm *ConfigReadModel) reduceConfigAddedEvent(e *ConfigAddedEvent) {
	rm.ConfigID = e.ConfigID
	rm.Name = e.Name
	rm.StylingType = e.StylingType
	rm.State = ConfigStateActive
}

func (rm *ConfigReadModel) reduceConfigChangedEvent(e *ConfigChangedEvent) {
	if e.Name != "" {
		rm.Name = e.Name
	}
	if e.StylingType.Valid() {
		rm.StylingType = e.StylingType
	}
}

func (rm *ConfigReadModel) reduceConfigStateChanged(configID string, state ConfigState) {
	rm.State = state
}
