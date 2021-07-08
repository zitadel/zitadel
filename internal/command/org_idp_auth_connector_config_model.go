package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/idpconfig"
	"github.com/caos/zitadel/internal/repository/org"
)

type OrgIDPAuthConnectorConfigWriteModel struct {
	AuthConnectorConfigWriteModel
}

func NewOrgIDPAuthConnectorConfigWriteModel(idpConfigID, orgID string) *OrgIDPAuthConnectorConfigWriteModel {
	return &OrgIDPAuthConnectorConfigWriteModel{
		AuthConnectorConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *OrgIDPAuthConnectorConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.IDPConfigAddedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.AuthConnectorConfigWriteModel.AppendEvents(e)
		case *org.IDPAuthConnectorConfigAddedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.AuthConnectorConfigWriteModel.AppendEvents(&e.AuthConnectorConfigAddedEvent)
		case *org.IDPAuthConnectorConfigChangedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.AuthConnectorConfigWriteModel.AppendEvents(&e.AuthConnectorConfigChangedEvent)
		case *org.IDPConfigReactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.AuthConnectorConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *org.IDPConfigDeactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.AuthConnectorConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *org.IDPConfigRemovedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.AuthConnectorConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		}
	}
}

func (wm *OrgIDPAuthConnectorConfigWriteModel) Reduce() error {
	return wm.AuthConnectorConfigWriteModel.Reduce()
}

func (wm *OrgIDPAuthConnectorConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.AggregateID).
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.IDPAuthConnectorConfigAddedEventType,
			org.IDPAuthConnectorConfigChangedEventType,
			org.IDPAuthConnectorMachineUserRemovedEventType,
			org.IDPConfigReactivatedEventType,
			org.IDPConfigDeactivatedEventType,
			org.IDPConfigRemovedEventType).
		Builder()
}

func (wm *OrgIDPAuthConnectorConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	baseURL,
	providerID,
	machineID string,
) (*org.IDPAuthConnectorConfigChangedEvent, bool, error) {

	changes := make([]idpconfig.AuthConnectorConfigChanges, 0)

	if wm.BaseURL != baseURL {
		changes = append(changes, idpconfig.ChangeBaseURL(baseURL))
	}
	if wm.ProviderID != providerID {
		changes = append(changes, idpconfig.ChangeProviderID(providerID))
	}
	if wm.MachineID != machineID {
		changes = append(changes, idpconfig.ChangeMachineID(machineID))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := org.NewIDPAuthConnectorConfigChangedEvent(ctx, aggregate, idpConfigID, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
