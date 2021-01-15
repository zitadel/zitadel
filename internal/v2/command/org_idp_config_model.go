package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idpconfig"
)

type OrgIDPConfigWriteModel struct {
	IDPConfigWriteModel
}

func NewOrgIDPConfigWriteModel(configID, orgID string) *OrgIDPConfigWriteModel {
	return &OrgIDPConfigWriteModel{
		IDPConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			ConfigID: configID,
		},
	}
}

func (wm *OrgIDPConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *OrgIDPConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.IDPConfigAddedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *iam.IDPConfigChangedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *iam.IDPConfigDeactivatedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *iam.IDPConfigReactivatedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *iam.IDPConfigRemovedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		case *iam.IDPOIDCConfigAddedEvent:
			if wm.ConfigID != e.IDPConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *iam.IDPOIDCConfigChangedEvent:
			if wm.ConfigID != e.IDPConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		}
	}
}

func (wm *OrgIDPConfigWriteModel) Reduce() error {
	return wm.IDPConfigWriteModel.Reduce()
}

func (wm *OrgIDPConfigWriteModel) AppendAndReduce(events ...eventstore.EventReader) error {
	wm.AppendEvents(events...)
	return wm.Reduce()
}

func (wm *OrgIDPConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	configID,
	name string,
	stylingType domain.IDPConfigStylingType,
) (*iam.IDPConfigChangedEvent, bool) {

	changes := make([]idpconfig.IDPConfigChanges, 0)
	if wm.Name != name {
		changes = append(changes, idpconfig.ChangeName(name))
	}
	if stylingType.Valid() && wm.StylingType != stylingType {
		changes = append(changes, idpconfig.ChangeStyleType(stylingType))
	}
	if len(changes) == 0 {
		return nil, false
	}
	return iam.NewIDPConfigChangedEvent(ctx, configID, changes), true
}
