package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type InstanceIDPOIDCConfigWriteModel struct {
	OIDCConfigWriteModel
}

func NewInstanceIDPOIDCConfigWriteModel(ctx context.Context, idpConfigID string) *InstanceIDPOIDCConfigWriteModel {
	return &InstanceIDPOIDCConfigWriteModel{
		OIDCConfigWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   authz.GetInstance(ctx).InstanceID(),
				ResourceOwner: authz.GetInstance(ctx).InstanceID(),
			},
			IDPConfigID: idpConfigID,
		},
	}
}

func (wm *InstanceIDPOIDCConfigWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.IDPOIDCConfigAddedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.OIDCConfigAddedEvent)
		case *instance.IDPOIDCConfigChangedEvent:
			if wm.IDPConfigID != e.IDPConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.OIDCConfigChangedEvent)
		case *instance.IDPConfigReactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.IDPConfigReactivatedEvent)
		case *instance.IDPConfigDeactivatedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.IDPConfigDeactivatedEvent)
		case *instance.IDPConfigRemovedEvent:
			if wm.IDPConfigID != e.ConfigID {
				continue
			}
			wm.OIDCConfigWriteModel.AppendEvents(&e.IDPConfigRemovedEvent)
		default:
			wm.OIDCConfigWriteModel.AppendEvents(e)
		}
	}
}

func (wm *InstanceIDPOIDCConfigWriteModel) Reduce() error {
	if err := wm.OIDCConfigWriteModel.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *InstanceIDPOIDCConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.IDPOIDCConfigAddedEventType,
			instance.IDPOIDCConfigChangedEventType,
			instance.IDPConfigReactivatedEventType,
			instance.IDPConfigDeactivatedEventType,
			instance.IDPConfigRemovedEventType).
		Builder()
}

func (wm *InstanceIDPOIDCConfigWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	clientID,
	issuer,
	authorizationEndpoint,
	tokenEndpoint,
	clientSecretString string,
	secretCrypto crypto.EncryptionAlgorithm,
	idpDisplayNameMapping,
	userNameMapping domain.OIDCMappingField,
	scopes ...string,
) (*instance.IDPOIDCConfigChangedEvent, bool, error) {

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
	if wm.AuthorizationEndpoint != authorizationEndpoint {
		changes = append(changes, idpconfig.ChangeAuthorizationEndpoint(authorizationEndpoint))
	}
	if wm.TokenEndpoint != tokenEndpoint {
		changes = append(changes, idpconfig.ChangeTokenEndpoint(tokenEndpoint))
	}
	if idpDisplayNameMapping.Valid() && wm.IDPDisplayNameMapping != idpDisplayNameMapping {
		changes = append(changes, idpconfig.ChangeIDPDisplayNameMapping(idpDisplayNameMapping))
	}
	if userNameMapping.Valid() && wm.UserNameMapping != userNameMapping {
		changes = append(changes, idpconfig.ChangeUserNameMapping(userNameMapping))
	}
	if !reflect.DeepEqual(wm.Scopes, scopes) {
		changes = append(changes, idpconfig.ChangeScopes(scopes))
	}
	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewIDPOIDCConfigChangedEvent(ctx, aggregate, idpConfigID, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
