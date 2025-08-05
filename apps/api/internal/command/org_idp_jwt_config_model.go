package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type IDPJWTConfigWriteModel struct {
	JWTConfigWriteModel
}

func NewOrgIDPJWTConfigWriteModel(idpConfigID, orgID string) *IDPJWTConfigWriteModel {
	return &IDPJWTConfigWriteModel{
		JWTConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *IDPJWTConfigWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *org.IDPJWTConfigAddedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *org.IDPJWTConfigChangedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *org.IDPConfigReactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *org.IDPConfigDeactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *org.IDPConfigRemovedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.JWTConfigWriteModel.AppendEvents(e)
		}
	}
}

func (wm *IDPJWTConfigWriteModel) Reduce() error {
	if err := wm.JWTConfigWriteModel.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *IDPJWTConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.IDPJWTConfigAddedEventType,
			org.IDPJWTConfigChangedEventType,
			org.IDPConfigReactivatedEventType,
			org.IDPConfigDeactivatedEventType,
			org.IDPConfigRemovedEventType).
		Builder()
}

func (wm *IDPJWTConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	jwtEndpoint,
	issuer,
	keysEndpoint,
	headerName string,
) (*org.IDPJWTConfigChangedEvent, bool, error) {

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
	changeEvent, err := org.NewIDPJWTConfigChangedEvent(ctx, aggregate, idpConfigID, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
