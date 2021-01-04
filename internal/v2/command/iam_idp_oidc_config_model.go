package command

import (
	"context"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"reflect"
)

type IDPOIDCConfigWriteModel struct {
	OIDCConfigWriteModel
}

func NewIDPOIDCConfigWriteModel(iamID, idpConfigID string) *IDPOIDCConfigWriteModel {
	return &IDPOIDCConfigWriteModel{
		OIDCConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID: iamID,
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *IDPOIDCConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.IDPOIDCConfigAddedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *iam.IDPOIDCConfigChangedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *iam.IDPConfigReactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *iam.IDPConfigDeactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *iam.IDPConfigRemovedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.OIDCConfigWriteModel.AppendEvents(e)
		}
	}
}

func (wm *IDPOIDCConfigWriteModel) Reduce() error {
	if err := wm.OIDCConfigWriteModel.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *IDPOIDCConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(wm.AggregateID)
}

func (wm *IDPOIDCConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	clientID,
	issuer,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	idpDisplayNameMapping,
	userNameMapping domain.OIDCMappingField,
	scopes ...string,
) (*iam.IDPOIDCConfigChangedEvent, bool, error) {
	hasChanged := false
	changedEvent := iam.NewIDPOIDCConfigChangedEvent(ctx)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, false, err
		}
		changedEvent.ClientSecret = clientSecret
	}
	if wm.ClientID != clientID {
		hasChanged = true
		changedEvent.ClientID = clientID
	}
	if wm.Issuer != issuer {
		hasChanged = true
		changedEvent.Issuer = issuer
	}
	if idpDisplayNameMapping.Valid() && wm.IDPDisplayNameMapping != idpDisplayNameMapping {
		hasChanged = true
		changedEvent.IDPDisplayNameMapping = idpDisplayNameMapping
	}
	if userNameMapping.Valid() && wm.UserNameMapping != userNameMapping {
		hasChanged = true
		changedEvent.UserNameMapping = userNameMapping
	}
	if reflect.DeepEqual(wm.Scopes, scopes) {
		hasChanged = true
		changedEvent.Scopes = scopes
	}
	return changedEvent, hasChanged, nil
}
