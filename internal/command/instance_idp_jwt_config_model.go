package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceIDPJWTConfigWriteModel struct {
	JWTConfigWriteModel
}

func NewInstanceIDPJWTConfigWriteModel(ctx context.Context, idpConfigID string) *InstanceIDPJWTConfigWriteModel {
	return &InstanceIDPJWTConfigWriteModel{
		JWTConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *InstanceIDPJWTConfigWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.IDPJWTConfigAddedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *instance.IDPJWTConfigChangedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *instance.IDPConfigReactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *instance.IDPConfigDeactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *instance.IDPConfigRemovedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.JWTConfigWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceIDPJWTConfigWriteModel) Reduce() error {
	if err := wm.JWTConfigWriteModel.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *InstanceIDPJWTConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.IDPJWTConfigAddedEventType,
			instance.IDPJWTConfigChangedEventType,
			instance.IDPConfigReactivatedEventType,
			instance.IDPConfigDeactivatedEventType,
			instance.IDPConfigRemovedEventType).
		Builder()
}

func (wm *InstanceIDPJWTConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	jwtEndpoint,
	issuer,
	keysEndpoint,
	headerName string,
) (*instance.IDPJWTConfigChangedEvent, bool, error) {

	changes := make([]idpconfig.JWTConfigChanges, 0)
	if wm.JWTEndpoint != jwtEndpoint {
		changes = append(changes, idpconfig.ChangeJWTEndpoint(jwtEndpoint))
	}
	if wm.Issuer != issuer {
		changes = append(changes, idpconfig.ChangeJWTIssuer(issuer))
	}
	if wm.KeysEndpoint != keysEndpoint {
		changes = append(changes, idpconfig.ChangeKeysEndpoint(keysEndpoint))
	}
	if wm.HeaderName != headerName {
		changes = append(changes, idpconfig.ChangeHeaderName(headerName))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewIDPJWTConfigChangedEvent(ctx, aggregate, idpConfigID, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
