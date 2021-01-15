package command

import (
	"context"
	"reflect"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idpconfig"
)

type IDPOIDCConfigWriteModel struct {
	OIDCConfigWriteModel
}

func NewOrgIDPOIDCConfigWriteModel(idpConfigID, orgID string) *IDPOIDCConfigWriteModel {
	return &IDPOIDCConfigWriteModel{
		OIDCConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   orgID,
				ResourceOwner: orgID,
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
		AggregateIDs(wm.AggregateID).
		ResourceOwner(wm.ResourceOwner)
}

func (wm *IDPOIDCConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	idpConfigID,
	clientID,
	issuer,
	clientSecretString string,
	secretCrypto crypto.Crypto,
	idpDisplayNameMapping,
	userNameMapping domain.OIDCMappingField,
	scopes ...string,
) (*iam.IDPOIDCConfigChangedEvent, bool, error) {

	changes := make([]idpconfig.OIDCConfigChanges, 0)
	var clientSecret *crypto.CryptoValue
	var err error
	if clientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(clientSecretString), secretCrypto)
		if err != nil {
			return nil, false, err
		}
		changes = append(changes, idpconfig.ChangeClientSecret(clientSecret))
	}
	if wm.ClientID != clientID {
		changes = append(changes, idpconfig.ChangeClientID(clientID))
	}
	if wm.Issuer != issuer {
		changes = append(changes, idpconfig.ChangeIssuer(issuer))
	}
	if idpDisplayNameMapping.Valid() && wm.IDPDisplayNameMapping != idpDisplayNameMapping {
		changes = append(changes, idpconfig.ChangeIDPDisplayNameMapping(idpDisplayNameMapping))
	}
	if userNameMapping.Valid() && wm.UserNameMapping != userNameMapping {
		changes = append(changes, idpconfig.ChangeUserNameMapping(userNameMapping))
	}
	if reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idpconfig.ChangeScopes(scopes))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	return iam.NewIDPOIDCConfigChangedEvent(ctx, idpConfigID, changes), true, nil
}
