package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
)

type IDPConfigWriteModel struct {
	eventstore.WriteModel

	State domain.IDPConfigState

	ConfigID     string
	Name         string
	AutoRegister bool
	StylingType  domain.IDPConfigStylingType
}

func (rm *IDPConfigWriteModel) AppendEvents(events ...eventstore.Event) {
	rm.WriteModel.AppendEvents(events...)
}

func (rm *IDPConfigWriteModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *idpconfig.IDPConfigAddedEvent:
			rm.reduceConfigAddedEvent(e)
		case *idpconfig.IDPConfigChangedEvent:
			rm.reduceConfigChangedEvent(e)
		case *idpconfig.IDPConfigDeactivatedEvent:
			rm.reduceConfigStateChanged(e.ConfigID, domain.IDPConfigStateInactive)
		case *idpconfig.IDPConfigReactivatedEvent:
			rm.reduceConfigStateChanged(e.ConfigID, domain.IDPConfigStateActive)
		case *idpconfig.IDPConfigRemovedEvent:
			rm.reduceConfigStateChanged(e.ConfigID, domain.IDPConfigStateRemoved)
		}
	}
	return rm.WriteModel.Reduce()
}

func (rm *IDPConfigWriteModel) reduceConfigAddedEvent(e *idpconfig.IDPConfigAddedEvent) {
	rm.ConfigID = e.ConfigID
	rm.Name = e.Name
	rm.StylingType = e.StylingType
	rm.AutoRegister = e.AutoRegister
	rm.State = domain.IDPConfigStateActive
}

func (rm *IDPConfigWriteModel) reduceConfigChangedEvent(e *idpconfig.IDPConfigChangedEvent) {
	if e.Name != nil {
		rm.Name = *e.Name
	}
	if e.StylingType != nil && e.StylingType.Valid() {
		rm.StylingType = *e.StylingType
	}
	if e.AutoRegister != nil {
		rm.AutoRegister = *e.AutoRegister
	}
}

func (rm *IDPConfigWriteModel) reduceConfigStateChanged(configID string, state domain.IDPConfigState) {
	rm.State = state
}

func (rm *IDPConfigWriteModel) Exists() bool {
	return rm.State.Exists()
}
