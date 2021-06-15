package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

type AuthConnectorConfigWriteModel struct {
	eventstore.WriteModel

	IDPConfigID        string
	BaseURL            string
	BackendConnectorID string
	State              domain.IDPConfigState
}

func (wm *AuthConnectorConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idpconfig.AuthConnectorConfigAddedEvent:
			wm.reduceConfigAddedEvent(e)
		case *idpconfig.AuthConnectorConfigChangedEvent:
			wm.reduceConfigChangedEvent(e)
		case *idpconfig.IDPConfigDeactivatedEvent:
			wm.State = domain.IDPConfigStateInactive
		case *idpconfig.IDPConfigReactivatedEvent:
			wm.State = domain.IDPConfigStateActive
		case *idpconfig.IDPConfigRemovedEvent:
			wm.State = domain.IDPConfigStateRemoved
		}
	}

	return wm.WriteModel.Reduce()
}

func (wm *AuthConnectorConfigWriteModel) reduceConfigAddedEvent(e *idpconfig.AuthConnectorConfigAddedEvent) {
	wm.IDPConfigID = e.IDPConfigID
	wm.BaseURL = e.BaseURL
	wm.BackendConnectorID = e.BackendConnectorID
	wm.State = domain.IDPConfigStateActive
}

func (wm *AuthConnectorConfigWriteModel) reduceConfigChangedEvent(e *idpconfig.AuthConnectorConfigChangedEvent) {
	if e.BaseURL != nil {
		wm.BaseURL = *e.BaseURL
	}
	if e.BackendConnectorID != nil {
		wm.BackendConnectorID = *e.BackendConnectorID
	}
}
