package projection

import (
	"context"
	"database/sql"
	"encoding/json"
	"slices"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceOAuthIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.OAuthIDPAddedEvent
	switch e := event.(type) {
	case *org.OAuthIDPAddedEvent:
		idpEvent = e.OAuthIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.OAuthIDPAddedEvent:
		idpEvent = e.OAuthIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yap9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPAddedEventType, instance.OAuthIDPAddedEventType})
	}

	oauth := domain.OAuth{
		ClientID:              idpEvent.ClientID,
		ClientSecret:          idpEvent.ClientSecret,
		AuthorizationEndpoint: idpEvent.AuthorizationEndpoint,
		TokenEndpoint:         idpEvent.TokenEndpoint,
		UserEndpoint:          idpEvent.UserEndpoint,
		Scopes:                idpEvent.Scopes,
		IDAttribute:           idpEvent.IDAttribute,
		UsePKCE:               idpEvent.UsePKCE,
	}

	payloadJSON, err := json.Marshal(oauth)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-mB2hq", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeOAuth),
			AllowCreation:     idpEvent.IsCreationAllowed,
			AllowLinking:      idpEvent.IsLinkingAllowed,
			AllowAutoCreation: idpEvent.IsAutoCreation,
			AllowAutoUpdate:   idpEvent.IsAutoUpdate,
			AutoLinkingField:  mapAutoLinkingField(idpEvent.AutoLinkingOption),
			Payload:           payloadJSON,
			CreatedAt:         event.CreatedAt(),
			UpdatedAt:         event.CreatedAt(),
		})
	}), nil
}

func (p *relationalTablesProjection) reduceOAuthIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.OAuthIDPChangedEvent
	switch e := event.(type) {
	case *org.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-K1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPChangedEventType, instance.OAuthIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		oauth, err := idpRepo.GetOAuth(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &oauth.OAuth
		payloadChanged := p.reduceOAuthIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			changes = append(changes, idpRepo.SetPayload(string(payloadJSON)))
		}

		changes = append(changes, idpRepo.SetUpdatedAt(gu.Ptr(event.CreatedAt())))
		_, err = idpRepo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(idpRepo, idpEvent.Aggregate().InstanceID, idpEvent.ID, orgId), changes...)
		return err

	}), nil
}

func (p *relationalTablesProjection) reduceOAuthIDPChangedColumns(payload *domain.OAuth, idpEvent *idp.OAuthIDPChangedEvent) bool {
	payloadChanged := false
	if payload == nil || idpEvent == nil {
		return payloadChanged
	}

	if idpEvent.ClientID != nil && *idpEvent.ClientID != payload.ClientID {
		payloadChanged = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChanged = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.AuthorizationEndpoint != nil && *idpEvent.AuthorizationEndpoint != payload.AuthorizationEndpoint {
		payloadChanged = true
		payload.AuthorizationEndpoint = *idpEvent.AuthorizationEndpoint
	}
	if idpEvent.TokenEndpoint != nil && *idpEvent.TokenEndpoint != payload.TokenEndpoint {
		payloadChanged = true
		payload.TokenEndpoint = *idpEvent.TokenEndpoint
	}
	if idpEvent.UserEndpoint != nil && *idpEvent.UserEndpoint != payload.UserEndpoint {
		payloadChanged = true
		payload.UserEndpoint = *idpEvent.UserEndpoint
	}
	if idpEvent.Scopes != nil && !slices.Equal(idpEvent.Scopes, payload.Scopes) {
		payloadChanged = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.IDAttribute != nil && *idpEvent.IDAttribute != payload.IDAttribute {
		payloadChanged = true
		payload.IDAttribute = *idpEvent.IDAttribute
	}
	if idpEvent.UsePKCE != nil && *idpEvent.UsePKCE != payload.UsePKCE {
		payloadChanged = true
		payload.UsePKCE = *idpEvent.UsePKCE
	}
	return payloadChanged
}
