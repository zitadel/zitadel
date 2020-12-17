package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/idpconfig"
)

type IDPConfigReadModel struct {
	eventstore.ReadModel

	State        domain.IDPConfigState
	ConfigID     string
	Name         string
	StylingType  domain.IDPConfigStylingType
	ProviderType domain.IdentityProviderType

	OIDCConfig *OIDCConfigReadModel
}

func NewIDPConfigReadModel(configID string) *IDPConfigReadModel {
	return &IDPConfigReadModel{
		ConfigID: configID,
	}
}

func (rm *IDPConfigReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *idpconfig.IDPConfigAddedEvent:
			rm.ReadModel.AppendEvents(e)
		case *idpconfig.IDPConfigChangedEvent:
			rm.ReadModel.AppendEvents(e)
		case *idpconfig.IDPConfigDeactivatedEvent:
			rm.ReadModel.AppendEvents(e)
		case *idpconfig.IDPConfigReactivatedEvent:
			rm.ReadModel.AppendEvents(e)
		case *idpconfig.IDPConfigRemovedEvent:
			rm.ReadModel.AppendEvents(e)
		case *idpconfig.OIDCConfigAddedEvent:
			rm.OIDCConfig = &OIDCConfigReadModel{}
			rm.ReadModel.AppendEvents(e)
			rm.OIDCConfig.AppendEvents(event)
		case *idpconfig.OIDCConfigChangedEvent:
			rm.ReadModel.AppendEvents(e)
			rm.OIDCConfig.AppendEvents(event)
		}
	}
}

func (rm *IDPConfigReadModel) Reduce() error {
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

	if rm.OIDCConfig != nil {
		if err := rm.OIDCConfig.Reduce(); err != nil {
			return err
		}
	}
	return rm.ReadModel.Reduce()
}

func (rm *IDPConfigReadModel) reduceConfigAddedEvent(e *idpconfig.IDPConfigAddedEvent) {
	rm.ConfigID = e.ConfigID
	rm.Name = e.Name
	rm.StylingType = e.StylingType
	rm.State = domain.IDPConfigStateActive
}

func (rm *IDPConfigReadModel) reduceConfigChangedEvent(e *idpconfig.IDPConfigChangedEvent) {
	if e.Name != "" {
		rm.Name = e.Name
	}
	if e.StylingType.Valid() {
		rm.StylingType = e.StylingType
	}
}

func (rm *IDPConfigReadModel) reduceConfigStateChanged(configID string, state domain.IDPConfigState) {
	rm.State = state
}
