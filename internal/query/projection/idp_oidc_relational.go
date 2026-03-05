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
	internal_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgID *string
	var idpEvent idpconfig.OIDCConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
		orgID = gu.Ptr(idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YFuAA", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
	}

	payloadJSON, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5cvzY", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		_, err = repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.IDPConfigID, orgID),
			repo.SetPayload(string(payloadJSON)),
			repo.SetType(domain.IDPTypeOIDC),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idpconfig.OIDCConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigChangedEventType, instance.IDPOIDCConfigChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-sh6Lp", "unable to cast to tx executer")
		}

		idpRepo := repository.IDProviderRepository()
		oidc, err := idpRepo.GetOIDC(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.IDPConfigID, orgId)))
		if err != nil {
			return err
		}

		if idpEvent.ClientID != nil && *idpEvent.ClientID != oidc.ClientID {
			oidc.ClientID = *idpEvent.ClientID
		}
		if idpEvent.ClientSecret != nil {
			oidc.ClientSecret = idpEvent.ClientSecret
		}
		if idpEvent.Issuer != nil && *idpEvent.Issuer != oidc.Issuer {
			oidc.Issuer = *idpEvent.Issuer
		}
		if idpEvent.AuthorizationEndpoint != nil && *idpEvent.AuthorizationEndpoint != oidc.AuthorizationEndpoint {
			oidc.AuthorizationEndpoint = *idpEvent.AuthorizationEndpoint
		}
		if idpEvent.TokenEndpoint != nil && *idpEvent.TokenEndpoint != oidc.TokenEndpoint {
			oidc.TokenEndpoint = *idpEvent.TokenEndpoint
		}
		if idpEvent.Scopes != nil && !slices.Equal(idpEvent.Scopes, oidc.Scopes) {
			oidc.Scopes = idpEvent.Scopes
		}
		if idpEvent.IDPDisplayNameMapping != nil {
			oidc.IDPDisplayNameMapping = mapOIDCMappingField(*idpEvent.IDPDisplayNameMapping)
		}
		if idpEvent.UserNameMapping != nil {
			oidc.UserNameMapping = mapOIDCMappingField(*idpEvent.UserNameMapping)
		}

		payloadJSON, err := json.Marshal(oidc.OIDC)
		if err != nil {
			return err
		}

		_, err = idpRepo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(idpRepo, idpEvent.Aggregate().InstanceID, idpEvent.IDPConfigID, orgId),
			idpRepo.SetPayload(string(payloadJSON)),
			idpRepo.SetType(domain.IDPTypeOIDC),
			idpRepo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil
}

func mapOIDCMappingField(field internal_domain.OIDCMappingField) domain.OIDCMappingField {
	switch field {
	case internal_domain.OIDCMappingFieldEmail:
		return domain.OIDCMappingFieldEmail
	case internal_domain.OIDCMappingFieldPreferredLoginName:
		return domain.OIDCMappingFieldPreferredLoginName
	case internal_domain.OIDCMappingFieldUnspecified:
		fallthrough
	default:
		return domain.OIDCMappingFieldUnspecified
	}
}

func (p *relationalTablesProjection) reduceOIDCIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.OIDCIDPAddedEvent
	switch e := event.(type) {
	case *org.OIDCIDPAddedEvent:
		idpEvent = e.OIDCIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.OIDCIDPAddedEvent:
		idpEvent = e.OIDCIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ys02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPAddedEventType, instance.OIDCIDPAddedEventType})
	}

	payloadJSON, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-C9ju3", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeOIDC),
			AllowCreation:     idpEvent.IsCreationAllowed,
			AllowAutoCreation: idpEvent.IsAutoCreation,
			AllowAutoUpdate:   idpEvent.IsAutoUpdate,
			AllowLinking:      idpEvent.IsLinkingAllowed,
			AutoLinkingField:  mapAutoLinkingField(idpEvent.AutoLinkingOption),
			Payload:           payloadJSON,
			CreatedAt:         event.CreatedAt(),
			UpdatedAt:         event.CreatedAt(),
		})
	}), nil
}

func (p *relationalTablesProjection) reduceOIDCIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.OIDCIDPChangedEvent
	switch e := event.(type) {
	case *org.OIDCIDPChangedEvent:
		idpEvent = e.OIDCIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.OIDCIDPChangedEvent:
		idpEvent = e.OIDCIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1K82ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-L8CQt", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		oidc, err := idpRepo.GetOIDC(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &oidc.OIDC
		payloadChanged := p.reduceOIDCIDPChangedColumns(payload, &idpEvent)
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

func (p *relationalTablesProjection) reduceOIDCIDPChangedColumns(payload *domain.OIDC, idpEvent *idp.OIDCIDPChangedEvent) bool {
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
	if idpEvent.Issuer != nil && *idpEvent.Issuer != payload.Issuer {
		payloadChanged = true
		payload.Issuer = *idpEvent.Issuer
	}
	if idpEvent.Scopes != nil && !slices.Equal(idpEvent.Scopes, payload.Scopes) {
		payloadChanged = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.IsIDTokenMapping != nil && *idpEvent.IsIDTokenMapping != payload.IsIDTokenMapping {
		payloadChanged = true
		payload.IsIDTokenMapping = *idpEvent.IsIDTokenMapping
	}
	if idpEvent.UsePKCE != nil && *idpEvent.UsePKCE != payload.UsePKCE {
		payloadChanged = true
		payload.UsePKCE = *idpEvent.UsePKCE
	}
	return payloadChanged
}
