package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
)

type IAMIDPConfigReadModel struct {
	IDPConfigReadModel

	iamID    string
	configID string
}

func NewIAMIDPConfigReadModel(iamID, configID string) *IAMIDPConfigReadModel {
	return &IAMIDPConfigReadModel{
		iamID:    iamID,
		configID: configID,
	}
}

func (rm *IAMIDPConfigReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.IDPConfigAddedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *iam.IDPConfigChangedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *iam.IDPConfigDeactivatedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *iam.IDPConfigReactivatedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *iam.IDPConfigRemovedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.IDPConfigRemovedEvent)
		case *iam.IDPOIDCConfigAddedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *iam.IDPOIDCConfigChangedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *iam.IDPAuthConnectorConfigAddedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.AuthConnectorConfigAddedEvent)
		case *iam.IDPAuthConnectorConfigChangedEvent:
			rm.IDPConfigReadModel.AppendEvents(&e.AuthConnectorConfigChangedEvent)
		}
	}
}

func (rm *IAMIDPConfigReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(rm.iamID).
		EventData(map[string]interface{}{
			"idpConfigId": rm.configID,
		}).Builder()
}
