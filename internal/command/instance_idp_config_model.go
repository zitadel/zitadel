package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceIDPConfigWriteModel struct {
	IDPConfigWriteModel
}

func NewInstanceIDPConfigWriteModel(ctx context.Context, configID string) *InstanceIDPConfigWriteModel {
	return &InstanceIDPConfigWriteModel{
		IDPConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
				InstanceID:    authz.GetInstance(ctx).InstanceID(),
			},
			ConfigID: configID,
		},
	}
}

func (wm *InstanceIDPConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.IDPConfigAddedEventType,
			instance.IDPConfigChangedEventType,
			instance.IDPConfigDeactivatedEventType,
			instance.IDPConfigReactivatedEventType,
			instance.IDPConfigRemovedEventType,
			instance.IDPOIDCConfigAddedEventType,
			instance.IDPOIDCConfigChangedEventType).
		Builder()
}

func (wm *InstanceIDPConfigWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.IDPConfigAddedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigAddedEvent)
		case *instance.IDPConfigChangedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigChangedEvent)
		case *instance.IDPConfigDeactivatedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *instance.IDPConfigReactivatedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *instance.IDPConfigRemovedEvent:
			if wm.ConfigID != e.ConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		case *instance.IDPOIDCConfigAddedEvent:
			if wm.ConfigID != e.IDPConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *instance.IDPOIDCConfigChangedEvent:
			if wm.ConfigID != e.IDPConfigID {
				continue
			}
			wm.IDPConfigWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		}
	}
}

func (wm *InstanceIDPConfigWriteModel) Reduce() error {
	return wm.IDPConfigWriteModel.Reduce()
}

func (wm *InstanceIDPConfigWriteModel) AppendAndReduce(events ...eventstore.Event) error {
	wm.AppendEvents(events...)
	return wm.Reduce()
}

func (wm *InstanceIDPConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	configID,
	name string,
	stylingType domain.IDPConfigStylingType,
	autoRegister bool,
) (*instance.IDPConfigChangedEvent, bool) {

	changes := make([]idpconfig.IDPConfigChanges, 0)
	oldName := ""
	if wm.Name != name {
		oldName = wm.Name
		changes = append(changes, idpconfig.ChangeName(name))
	}
	if stylingType.Valid() && wm.StylingType != stylingType {
		changes = append(changes, idpconfig.ChangeStyleType(stylingType))
	}
	if wm.AutoRegister != autoRegister {
		changes = append(changes, idpconfig.ChangeAutoRegister(autoRegister))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changeEvent, err := instance.NewIDPConfigChangedEvent(ctx, aggregate, configID, oldName, changes)
	if err != nil {
		return nil, false
	}
	return changeEvent, true
}
