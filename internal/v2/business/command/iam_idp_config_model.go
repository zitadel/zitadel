package command

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

type IAMIDPConfigWriteModel struct {
	IDPConfigWriteModel
}

func NewIAMIDPConfigWriteModel(iamID, configID string) *IAMIDPConfigWriteModel {
	return &IAMIDPConfigWriteModel{
		IDPConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
			ConfigID: configID,
		},
	}
}

func (wm *IAMIDPConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func (wm *IAMIDPConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
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

func (wm *IAMIDPConfigWriteModel) Reduce() error {
	return wm.IDPConfigWriteModel.Reduce()
}

func (wm *IAMIDPConfigWriteModel) AppendAndReduce(events ...eventstore.EventReader) error {
	wm.AppendEvents(events...)
	return wm.Reduce()
}

func (wm *IAMIDPConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	configID,
	name string,
	stylingType domain.IDPConfigStylingType,
) (*iam.IDPConfigChangedEvent, bool) {

	hasChanged := false
	changedEvent := iam.NewIDPConfigChangedEvent(ctx)
	changedEvent.ConfigID = configID
	if wm.Name != name {
		hasChanged = true
		changedEvent.Name = name
	}
	if stylingType.Valid() && wm.StylingType != stylingType {
		hasChanged = true
		changedEvent.StylingType = stylingType
	}
	return changedEvent, hasChanged
}
