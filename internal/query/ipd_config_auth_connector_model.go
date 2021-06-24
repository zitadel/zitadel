package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

type IDPAuthConnectorConfigReadModel struct {
	eventstore.ReadModel

	IDPConfigID string
	BaseURL     string
	ProviderID  string
	MachineID   string
}

func (rm *IDPAuthConnectorConfigReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *idpconfig.AuthConnectorConfigAddedEvent:
			rm.reduceConfigAddedEvent(e)
		case *idpconfig.AuthConnectorConfigChangedEvent:
			rm.reduceConfigChangedEvent(e)
		}
	}

	return rm.ReadModel.Reduce()
}

func (rm *IDPAuthConnectorConfigReadModel) reduceConfigAddedEvent(e *idpconfig.AuthConnectorConfigAddedEvent) {
	rm.IDPConfigID = e.IDPConfigID
	rm.BaseURL = e.BaseURL
	rm.ProviderID = e.ProviderID
	rm.MachineID = e.MachineID
}

func (rm *IDPAuthConnectorConfigReadModel) reduceConfigChangedEvent(e *idpconfig.AuthConnectorConfigChangedEvent) {
	if e.BaseURL != nil {
		rm.BaseURL = *e.BaseURL
	}
	if e.ProviderID != nil {
		rm.ProviderID = *e.ProviderID
	}
	if e.MachineID != nil {
		rm.MachineID = *e.MachineID
	}
}
