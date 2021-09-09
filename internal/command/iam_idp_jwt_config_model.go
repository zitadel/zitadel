package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

type IAMIDPJWTConfigWriteModel struct {
	JWTConfigWriteModel
}

func NewIAMIDPJWTConfigWriteModel(idpConfigID string) *IAMIDPJWTConfigWriteModel {
	return &IAMIDPJWTConfigWriteModel{
		JWTConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   domain.IAMID,
				ResourceOwner: domain.IAMID,
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *IAMIDPJWTConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.IDPJWTConfigAddedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.JWTConfigAddedEvent)
		case *iam.IDPJWTConfigChangedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.JWTConfigChangedEvent)
		case *iam.IDPConfigReactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *iam.IDPConfigDeactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *iam.IDPConfigRemovedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.JWTConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.JWTConfigWriteModel.AppendEvents(e)
		}
	}
}

func (wm *IAMIDPJWTConfigWriteModel) Reduce() error {
	if err := wm.JWTConfigWriteModel.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *IAMIDPJWTConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(iam.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			iam.IDPJWTConfigAddedEventType,
			iam.IDPJWTConfigChangedEventType,
			iam.IDPConfigReactivatedEventType,
			iam.IDPConfigDeactivatedEventType,
			iam.IDPConfigRemovedEventType).
		Builder()
}

func (wm *IAMIDPJWTConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	issuer,
	keysEndpoint string,
) (*iam.IDPJWTConfigChangedEvent, bool, error) {

	changes := make([]idpconfig.JWTConfigChanges, 0)
	if wm.Issuer != issuer {
		changes = append(changes, idpconfig.ChangeJWTIssuer(issuer))
	}
	if wm.KeysEndpoint != keysEndpoint {
		changes = append(changes, idpconfig.ChangeKeysEndpoint(keysEndpoint))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := iam.NewIDPJWTConfigChangedEvent(ctx, aggregate, idpConfigID, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
