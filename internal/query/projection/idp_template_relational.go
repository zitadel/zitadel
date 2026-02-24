package projection

import (
	"context"
	"database/sql"
	"encoding/json"

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

const (
	IDPOrgIdCol             = "org_id"
	IDPPayloadCol           = "payload"
	IDPOrgId                = "org_id"
	IDPAllowCreationCol     = "allow_creation"
	IDPAllowLinkingCol      = "allow_linking"
	IDPAllowAutoCreationCol = "allow_auto_creation"
	IDPAllowAutoUpdateCol   = "allow_auto_update"
	IDPAllowAutoLinkingCol  = "auto_linking_field"
)

func (p *relationalTablesProjection) reduceIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idpconfig.IDPConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.IDPConfigAddedEvent:
		idpEvent = e.IDPConfigAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YcUdQ", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigAddedEventType, instance.IDPConfigAddedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:   event.Aggregate().InstanceID,
			OrgID:        orgId,
			ID:           idpEvent.ConfigID,
			State:        domain.IDPStateActive,
			Name:         idpEvent.Name,
			Type:         mapIDPConfigType(idpEvent.Typ),
			AutoRegister: idpEvent.AutoRegister,
		})

	}), nil
}

func mapIDPConfigType(typ internal_domain.IDPConfigType) *domain.IDPType {
	switch typ {
	case internal_domain.IDPConfigTypeOIDC:
		return gu.Ptr(domain.IDPTypeOIDC)
	case internal_domain.IDPConfigTypeSAML:
		return gu.Ptr(domain.IDPTypeSAML)
	case internal_domain.IDPConfigTypeJWT:
		return gu.Ptr(domain.IDPTypeJWT)
	default:
		return nil
	}
}

func (p *relationalTablesProjection) reduceIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgCond handler.Condition
	var idpEvent idpconfig.IDPConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
		orgCond = handler.NewCond(IDPOrgId, idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YVvJD", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
	}

	repo := repository.IDProviderRepository()

	changes := make(database.Changes, 0, 3)
	if idpEvent.Name != nil {
		changes = append(changes, repo.SetName(*idpEvent.Name))
	}
	if idpEvent.AutoRegister != nil {
		changes = append(changes, repo.SetAutoRegister(*idpEvent.AutoRegister))
	}
	if len(changes) == 0 {
		return handler.NewNoOpStatement(&idpEvent), nil
	}
	changes = append(changes, repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())))

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-9sX8h", "reduce.wrong.db.pool %T", ex)
		}

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx), repo.PrimaryKeyCondition(event.Aggregate().InstanceID, idpEvent.ConfigID), changes...)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceIDPDeactivated(event eventstore.Event) (*handler.Statement, error) {
	var orgCond handler.Condition
	var idpEvent idpconfig.IDPConfigDeactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
		orgCond = handler.NewCond(IDPOrgId, idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y4O5l", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigDeactivatedEventType, instance.IDPConfigDeactivatedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-9sX8h", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx), repo.PrimaryKeyCondition(event.Aggregate().InstanceID, idpEvent.ConfigID),
			repo.SetState(domain.IDPStateInactive),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPStateCol, domain.IDPStateInactive),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			orgCond,
		},
	), nil
}

func (p *relationalTablesProjection) reduceIDPReactivated(event eventstore.Event) (*handler.Statement, error) {
	var orgCond handler.Condition
	var idpEvent idpconfig.IDPConfigReactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
		orgCond = handler.NewCond(IDPOrgId, idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y8QyS", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigReactivatedEventType, instance.IDPConfigReactivatedEventType})
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPStateCol, domain.IDPStateActive.String()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			orgCond,
		},
	), nil
}

// func (p *relationalTablesProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
// 	var orgCond handler.Condition
// 	var idpEvent idpconfig.IDPConfigRemovedEvent
// 	switch e := event.(type) {
// 	case *org.IDPConfigRemovedEvent:
// 		idpEvent = e.IDPConfigRemovedEvent
// 		orgCond = handler.NewCond(IDPOrgId, idpEvent.Aggregate().ResourceOwner)
// 	case *instance.IDPConfigRemovedEvent:
// 		idpEvent = e.IDPConfigRemovedEvent
// 		orgCond = handler.NewIsNullCond((IDPOrgId))
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y4zy8", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
// 	}

// 	return handler.NewDeleteStatement(
// 		&idpEvent,
// 		[]handler.Condition{
// 			handler.NewCond(IDPIDCol, idpEvent.ConfigID),
// 			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
// 			orgCond,
// 		},
// 	), nil
// }

func (p *relationalTablesProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
	var (
		orgCond handler.Condition
		idpID   string
	)
	switch e := event.(type) {
	case *org.IDPRemovedEvent:
		idpID = e.ID
		orgCond = handler.NewCond(IDPOrgId, e.RemovedEvent.Aggregate().ResourceOwner)
	case *instance.IDPRemovedEvent:
		idpID = e.ID
		orgCond = handler.NewIsNullCond((IDPOrgId))
	case *org.IDPConfigRemovedEvent:
		idpID = e.ID
		orgCond = handler.NewCond(IDPOrgId, e.IDPConfigRemovedEvent.Aggregate().ResourceOwner)
	case *instance.IDPConfigRemovedEvent:
		idpID = e.ID
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ybcvwin2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPRemovedEventType, instance.IDPRemovedEventType})
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(IDPTemplateIDCol, idpID),
			handler.NewCond(IDPTemplateInstanceIDCol, event.Aggregate().InstanceID),
			orgCond,
		},
	), nil
}

func (p *relationalTablesProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgCond handler.Condition
	var idpEvent idpconfig.OIDCConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
		orgCond = handler.NewCond(IDPOrgId, idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPOIDCConfigAddedEvent:
		idpEvent = e.OIDCConfigAddedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YFuAA", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigAddedEventType, instance.IDPOIDCConfigAddedEventType})
	}

	payloadJSON, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(IDPTypeCol, domain.IDPTypeOIDC),
			handler.NewCol(UpdatedAt, idpEvent.CreatedAt()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			orgCond,
		},
	), nil
}

func (p *relationalTablesProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idpconfig.OIDCConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.IDPOIDCConfigChangedEvent:
		idpEvent = e.OIDCConfigChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPOIDCConfigChangedEventType, instance.IDPOIDCConfigChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-sh6Lp", "unable to cast to tx executer")
		}

		idpRepo := repository.IDProviderRepository()
		oidc, err := idpRepo.GetOIDC(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.IDPConfigID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		if idpEvent.ClientID != nil {
			oidc.ClientID = *idpEvent.ClientID
		}
		if idpEvent.ClientSecret != nil {
			oidc.ClientSecret = idpEvent.ClientSecret
		}
		if idpEvent.Issuer != nil {
			oidc.Issuer = *idpEvent.Issuer
		}
		if idpEvent.AuthorizationEndpoint != nil {
			oidc.AuthorizationEndpoint = *idpEvent.AuthorizationEndpoint
		}
		if idpEvent.TokenEndpoint != nil {
			oidc.TokenEndpoint = *idpEvent.TokenEndpoint
		}
		if idpEvent.Scopes != nil {
			oidc.Scopes = idpEvent.Scopes
		}
		if idpEvent.IDPDisplayNameMapping != nil {
			oidc.IDPDisplayNameMapping = domain.OIDCMappingField(*idpEvent.IDPDisplayNameMapping)
		}
		if idpEvent.UserNameMapping != nil {
			oidc.UserNameMapping = domain.OIDCMappingField(*idpEvent.UserNameMapping)
		}

		payloadJSON, err := json.Marshal(idpEvent)
		if err != nil {
			return err
		}

		return handler.NewUpdateStatement(
			&idpEvent,
			[]handler.Column{
				handler.NewCol(IDPPayloadCol, payloadJSON),
				handler.NewCol(IDPTypeCol, domain.IDPTypeOIDC),
				handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
			},
			[]handler.Condition{
				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
				handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)
	}), nil
}

func (p *relationalTablesProjection) reduceJWTConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgCond handler.Condition
	var idpEvent idpconfig.JWTConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
		orgCond = handler.NewCond(IDPOrgId, idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YvPdb", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
	}

	payloadJSON, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(IDPTypeCol, domain.IDPTypeJWT),
			handler.NewCol(UpdatedAt, idpEvent.CreatedAt()),
		},
		[]handler.Condition{
			handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
			handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
			orgCond,
		},
	), nil
}

func (p *relationalTablesProjection) reduceJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idpconfig.JWTConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigChangedEventType, instance.IDPJWTConfigChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-sh6Lp", "unable to cast to tx executer")
		}

		idpRepo := repository.IDProviderRepository()
		jwt, err := idpRepo.GetJWT(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.IDPConfigID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		if idpEvent.JWTEndpoint != nil {
			jwt.JWTEndpoint = *idpEvent.JWTEndpoint
		}
		if idpEvent.Issuer != nil {
			jwt.Issuer = *idpEvent.Issuer
		}
		if idpEvent.KeysEndpoint != nil {
			jwt.KeysEndpoint = *idpEvent.KeysEndpoint
		}
		if idpEvent.HeaderName != nil {
			jwt.HeaderName = *idpEvent.HeaderName
		}

		payloadJSON, err := json.Marshal(idpEvent)
		if err != nil {
			return err
		}

		return handler.NewUpdateStatement(
			&idpEvent,
			[]handler.Column{
				handler.NewCol(IDPPayloadCol, payloadJSON),
				handler.NewCol(IDPTypeCol, domain.IDPTypeJWT),
				handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
			},
			[]handler.Condition{
				handler.NewCond(IDPIDCol, idpEvent.IDPConfigID),
				handler.NewCond(IDPInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)
	}), nil
}

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

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOAuth),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceOAuthIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.OAuthIDPChangedEvent
	switch e := event.(type) {
	case *org.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.OAuthIDPChangedEvent:
		idpEvent = e.OAuthIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-K1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OAuthIDPChangedEventType, instance.OAuthIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		oauth, err := idpRepo.GetOAuth(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &oauth.OAuth
		payloadChanged := p.reduceOAuthIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
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

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeOIDC),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceOIDCIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.OIDCIDPChangedEvent
	switch e := event.(type) {
	case *org.OIDCIDPChangedEvent:
		idpEvent = e.OIDCIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.OIDCIDPChangedEvent:
		idpEvent = e.OIDCIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1K82ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPChangedEventType, instance.OIDCIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-L8CQt", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		oidc, err := idpRepo.GetOIDC(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &oidc.OIDC
		payloadChanged := p.reduceOIDCIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)
	}), nil
}

func (p *relationalTablesProjection) reduceOIDCIDPMigratedAzureAD(event eventstore.Event) (*handler.Statement, error) {
	var orgCond handler.Condition
	var idpEvent idp.OIDCIDPMigratedAzureADEvent
	switch e := event.(type) {
	case *org.OIDCIDPMigratedAzureADEvent:
		idpEvent = e.OIDCIDPMigratedAzureADEvent
		orgCond = handler.NewCond(IDPOrgId, idpEvent.Aggregate().ResourceOwner)
	case *instance.OIDCIDPMigratedAzureADEvent:
		idpEvent = e.OIDCIDPMigratedAzureADEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yb582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedAzureADEventType, instance.OIDCIDPMigratedAzureADEventType})
	}

	azureTenant, err := domain.AzureTenantTypeString(idpEvent.Tenant)
	if err != nil {
		return nil, err
	}

	azure := domain.Azure{
		ClientID:        idpEvent.ClientID,
		ClientSecret:    idpEvent.ClientSecret,
		Scopes:          idpEvent.Scopes,
		Tenant:          azureTenant,
		IsEmailVerified: idpEvent.IsEmailVerified,
	}

	payloadJSON, err := json.Marshal(azure)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeAzure),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			orgCond,
		},
	), nil
}

func (p *relationalTablesProjection) reduceOIDCIDPMigratedGoogle(event eventstore.Event) (*handler.Statement, error) {
	var orgCond handler.Condition
	var idpEvent idp.OIDCIDPMigratedGoogleEvent
	switch e := event.(type) {
	case *org.OIDCIDPMigratedGoogleEvent:
		idpEvent = e.OIDCIDPMigratedGoogleEvent
		orgCond = handler.NewCond(IDPOrgId, idpEvent.Aggregate().ResourceOwner)
	case *instance.OIDCIDPMigratedGoogleEvent:
		idpEvent = e.OIDCIDPMigratedGoogleEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y1502hk", "reduce.wrong.event.type %v", []eventstore.EventType{org.OIDCIDPMigratedGoogleEventType, instance.OIDCIDPMigratedGoogleEventType})
	}

	google := domain.Google{
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payloadJSON, err := json.Marshal(google)
	if err != nil {
		return nil, err
	}

	return handler.NewUpdateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGoogle),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			orgCond,
		},
	), nil
}

func (p *relationalTablesProjection) reduceJWTIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.JWTIDPAddedEvent
	switch e := event.(type) {
	case *org.JWTIDPAddedEvent:
		idpEvent = e.JWTIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.JWTIDPAddedEvent:
		idpEvent = e.JWTIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yopi2s", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPAddedEventType, instance.JWTIDPAddedEventType})
	}

	jwt := domain.JWT{
		JWTEndpoint:  idpEvent.JWTEndpoint,
		Issuer:       idpEvent.Issuer,
		KeysEndpoint: idpEvent.KeysEndpoint,
		HeaderName:   idpEvent.HeaderName,
	}

	payloadJSON, err := json.Marshal(jwt)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeJWT),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceJWTIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.JWTIDPChangedEvent
	switch e := event.(type) {
	case *org.JWTIDPChangedEvent:
		idpEvent = e.JWTIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.JWTIDPChangedEvent:
		idpEvent = e.JWTIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-H15j2il", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		jwt, err := idpRepo.GetJWT(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &jwt.JWT
		payloadChanged := p.reduceJWTIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)
	}), nil
}

func (p *relationalTablesProjection) reduceAzureADIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.AzureADIDPAddedEvent
	switch e := event.(type) {
	case *org.AzureADIDPAddedEvent:
		idpEvent = e.AzureADIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.AzureADIDPAddedEvent:
		idpEvent = e.AzureADIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPAddedEventType, instance.AzureADIDPAddedEventType})
	}

	azureTenant, err := domain.AzureTenantTypeString(idpEvent.Tenant)
	if err != nil {
		return nil, err
	}

	azure := domain.Azure{
		ClientID:        idpEvent.ClientID,
		ClientSecret:    idpEvent.ClientSecret,
		Scopes:          idpEvent.Scopes,
		Tenant:          azureTenant,
		IsEmailVerified: idpEvent.IsEmailVerified,
	}

	payloadJSON, err := json.Marshal(azure)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeAzure),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceAzureADIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.AzureADIDPChangedEvent
	switch e := event.(type) {
	case *org.AzureADIDPChangedEvent:
		idpEvent = e.AzureADIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.AzureADIDPChangedEvent:
		idpEvent = e.AzureADIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YZ5x25s", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPChangedEventType, instance.AzureADIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		azure, err := idpRepo.GetAzureAD(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &azure.Azure
		payloadChanged, err := p.reduceAzureADIDPChangedColumns(payload, &idpEvent)
		if err != nil {
			return err
		}

		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceGitHubIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GitHubIDPAddedEvent
	switch e := event.(type) {
	case *org.GitHubIDPAddedEvent:
		idpEvent = e.GitHubIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GitHubIDPAddedEvent:
		idpEvent = e.GitHubIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPAddedEventType, instance.GitHubIDPAddedEventType})
	}

	github := domain.Github{
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payloadJSON, err := json.Marshal(github)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitHub),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceGitHubIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.GitHubIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-L1U89ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPChangedEventType, instance.GitHubIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		github, err := idpRepo.GetGithub(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &github.Github
		payloadChanged := p.reduceGitHubIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceGitHubEnterpriseIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GitHubEnterpriseIDPAddedEvent
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPAddedEvent:
		idpEvent = e.GitHubEnterpriseIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GitHubEnterpriseIDPAddedEvent:
		idpEvent = e.GitHubEnterpriseIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yf3g2a", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPAddedEventType, instance.GitHubEnterpriseIDPAddedEventType})
	}

	githubEnterprise := domain.GithubEnterprise{
		ClientID:              idpEvent.ClientID,
		ClientSecret:          idpEvent.ClientSecret,
		AuthorizationEndpoint: idpEvent.AuthorizationEndpoint,
		TokenEndpoint:         idpEvent.TokenEndpoint,
		UserEndpoint:          idpEvent.UserEndpoint,
		Scopes:                idpEvent.Scopes,
	}

	payloadJSON, err := json.Marshal(githubEnterprise)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitHubEnterprise),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceGitHubEnterpriseIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.GitHubEnterpriseIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YDg3g", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPChangedEventType, instance.GitHubEnterpriseIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		githubEnterprise, err := idpRepo.GetGithubEnterprise(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &githubEnterprise.GithubEnterprise
		payloadChanged := p.reduceGitHubEnterpriseIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceGitLabIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GitLabIDPAddedEvent
	switch e := event.(type) {
	case *org.GitLabIDPAddedEvent:
		idpEvent = e.GitLabIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GitLabIDPAddedEvent:
		idpEvent = e.GitLabIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y9a022b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPAddedEventType, instance.GitLabIDPAddedEventType})
	}

	gitlab := domain.Gitlab{
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payloadJSON, err := json.Marshal(gitlab)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitLab),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceGitLabIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.GitLabIDPChangedEvent
	switch e := event.(type) {
	case *org.GitLabIDPChangedEvent:
		idpEvent = e.GitLabIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.GitLabIDPChangedEvent:
		idpEvent = e.GitLabIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-mT5827b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPChangedEventType, instance.GitLabIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		gitlab, err := idpRepo.GetGitlab(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &gitlab.Gitlab
		payloadChanged := p.reduceGitLabIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceGitLabSelfHostedIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GitLabSelfHostedIDPAddedEvent
	switch e := event.(type) {
	case *org.GitLabSelfHostedIDPAddedEvent:
		idpEvent = e.GitLabSelfHostedIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GitLabSelfHostedIDPAddedEvent:
		idpEvent = e.GitLabSelfHostedIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YAF3gw", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPAddedEventType, instance.GitLabSelfHostedIDPAddedEventType})
	}

	gitlabSelfHosting := domain.GitlabSelfHosting{
		Issuer:       idpEvent.Issuer,
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payloadJSON, err := json.Marshal(gitlabSelfHosting)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGitLabSelfHosted),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceGitLabSelfHostedIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.GitLabSelfHostedIDPChangedEvent
	switch e := event.(type) {
	case *org.GitLabSelfHostedIDPChangedEvent:
		idpEvent = e.GitLabSelfHostedIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.GitLabSelfHostedIDPChangedEvent:
		idpEvent = e.GitLabSelfHostedIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YAf3g2", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPChangedEventType, instance.GitLabSelfHostedIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		gitlabSelfHosted, err := idpRepo.GetGitlabSelfHosting(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &gitlabSelfHosted.GitlabSelfHosting
		payloadChanged := p.reduceGitLabSelfHostedIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceGoogleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GoogleIDPAddedEvent
	switch e := event.(type) {
	case *org.GoogleIDPAddedEvent:
		idpEvent = e.GoogleIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GoogleIDPAddedEvent:
		idpEvent = e.GoogleIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Yp9ihb", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPAddedEventType, instance.GoogleIDPAddedEventType})
	}

	google := domain.Google{
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payloadJSON, err := json.Marshal(google)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeGoogle),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceGoogleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.GoogleIDPChangedEvent
	switch e := event.(type) {
	case *org.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YN58hml", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPChangedEventType, instance.GoogleIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		google, err := idpRepo.GetGoogle(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &google.Google
		payloadChanged := p.reduceGoogleIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceLDAPIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.LDAPIDPAddedEvent
	switch e := event.(type) {
	case *org.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.LDAPIDPAddedEvent:
		idpEvent = e.LDAPIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9s02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPAddedEventType, instance.LDAPIDPAddedEventType})
	}

	ldap := domain.LDAP{
		Servers:           idpEvent.Servers,
		StartTLS:          idpEvent.StartTLS,
		BaseDN:            idpEvent.BaseDN,
		BindDN:            idpEvent.BindDN,
		BindPassword:      idpEvent.BindPassword,
		UserBase:          idpEvent.UserBase,
		UserObjectClasses: idpEvent.UserObjectClasses,
		UserFilters:       idpEvent.UserFilters,
		Timeout:           idpEvent.Timeout,
		LDAPAttributes: domain.LDAPAttributes{
			IDAttribute:                idpEvent.IDAttribute,
			FirstNameAttribute:         idpEvent.FirstNameAttribute,
			LastNameAttribute:          idpEvent.LastNameAttribute,
			DisplayNameAttribute:       idpEvent.DisplayNameAttribute,
			NickNameAttribute:          idpEvent.NickNameAttribute,
			PreferredUsernameAttribute: idpEvent.PreferredUsernameAttribute,
			EmailAttribute:             idpEvent.EmailAttribute,
			EmailVerifiedAttribute:     idpEvent.EmailVerifiedAttribute,
			PhoneAttribute:             idpEvent.PhoneAttribute,
			PhoneVerifiedAttribute:     idpEvent.PhoneVerifiedAttribute,
			PreferredLanguageAttribute: idpEvent.PreferredLanguageAttribute,
			AvatarURLAttribute:         idpEvent.AvatarURLAttribute,
			ProfileAttribute:           idpEvent.ProfileAttribute,
		},
	}

	payloadJSON, err := json.Marshal(ldap)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeLDAP),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceLDAPIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.LDAPIDPChangedEvent
	switch e := event.(type) {
	case *org.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPChangedEventType, instance.LDAPIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		ldap, err := idpRepo.GetLDAP(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &ldap.LDAP
		payloadChanged := p.reduceLDAPIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceAppleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.AppleIDPAddedEvent
	switch e := event.(type) {
	case *org.AppleIDPAddedEvent:
		idpEvent = e.AppleIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.AppleIDPAddedEvent:
		idpEvent = e.AppleIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YFvg3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPAddedEventType /*, instance.AppleIDPAddedEventType*/})
	}

	apple := domain.Apple{
		ClientID:   idpEvent.ClientID,
		TeamID:     idpEvent.TeamID,
		KeyID:      idpEvent.KeyID,
		PrivateKey: idpEvent.PrivateKey,
		Scopes:     idpEvent.Scopes,
	}

	payloadJSON, err := json.Marshal(apple)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeApple),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceAppleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.AppleIDPChangedEvent
	switch e := event.(type) {
	case *org.AppleIDPChangedEvent:
		idpEvent = e.AppleIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.AppleIDPChangedEvent:
		idpEvent = e.AppleIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YBez3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPChangedEventType /*, instance.AppleIDPChangedEventType*/})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		apple, err := idpRepo.GetApple(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &apple.Apple
		payloadChanged := p.reduceAppleIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceSAMLIDPAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.SAMLIDPAddedEvent
	switch e := event.(type) {
	case *org.SAMLIDPAddedEvent:
		idpEvent = e.SAMLIDPAddedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.SAMLIDPAddedEvent:
		idpEvent = e.SAMLIDPAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ys02m1", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPAddedEventType, instance.SAMLIDPAddedEventType})
	}

	saml := domain.SAML{
		Metadata:                      idpEvent.Metadata,
		Key:                           idpEvent.Key,
		Certificate:                   idpEvent.Certificate,
		Binding:                       idpEvent.Binding,
		WithSignedRequest:             idpEvent.WithSignedRequest,
		NameIDFormat:                  idpEvent.NameIDFormat,
		TransientMappingAttributeName: idpEvent.TransientMappingAttributeName,
		SignatureAlgorithm:            idpEvent.SignatureAlgorithm,
	}

	payloadJSON, err := json.Marshal(saml)
	if err != nil {
		return nil, err
	}

	return handler.NewCreateStatement(
		&idpEvent,
		[]handler.Column{
			handler.NewCol(IDPTemplateIDCol, idpEvent.ID),
			handler.NewCol(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
			handler.NewCol(IDPOrgId, orgId),
			handler.NewCol(IDPTemplateNameCol, idpEvent.Name),
			handler.NewCol(IDPTemplateTypeCol, domain.IDPTypeSAML),
			handler.NewCol(IDPTemplateStateCol, domain.IDPStateActive.String()),
			handler.NewCol(IDPAllowCreationCol, idpEvent.IsCreationAllowed),
			handler.NewCol(IDPAllowLinkingCol, idpEvent.IsLinkingAllowed),
			handler.NewCol(IDPAllowAutoCreationCol, idpEvent.IsAutoCreation),
			handler.NewCol(IDPAllowAutoUpdateCol, idpEvent.IsAutoUpdate),
			handler.NewCol(IDPAllowAutoLinkingCol, func() any {
				if idpEvent.AutoLinkingOption == internal_domain.AutoLinkingOptionUnspecified {
					return nil
				}
				return domain.IDPAutoLinkingField(idpEvent.AutoLinkingOption)
			}()),
			handler.NewCol(IDPPayloadCol, payloadJSON),
			handler.NewCol(CreatedAt, idpEvent.CreationDate()),
			handler.NewCol(UpdatedAt, idpEvent.CreationDate()),
		},
	), nil
}

func (p *relationalTablesProjection) reduceSAMLIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var orgCond handler.Condition
	var idpEvent idp.SAMLIDPChangedEvent
	switch e := event.(type) {
	case *org.SAMLIDPChangedEvent:
		idpEvent = e.SAMLIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
		orgCond = handler.NewCond(IDPOrgId, orgId)
	case *instance.SAMLIDPChangedEvent:
		idpEvent = e.SAMLIDPChangedEvent
		orgCond = handler.NewIsNullCond((IDPOrgId))
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y7c0fii4ad", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPChangedEventType, instance.SAMLIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		saml, err := idpRepo.GetSAML(ctx, v3_sql.SQLTx(tx), idpRepo.IDCondition(idpEvent.ID), idpEvent.Agg.InstanceID, orgId)
		if err != nil {
			return err
		}

		columns := p.reduceIDPChangedTemplateColumns(idpEvent.Name, idpEvent.OptionChanges)

		payload := &saml.SAML
		payloadChanged := p.reduceSAMLIDPChangedColumns(payload, &idpEvent)
		if payloadChanged {
			payloadJSON, err := json.Marshal(payload)
			if err != nil {
				return err
			}
			columns = append(columns, handler.NewCol(IDPPayloadCol, payloadJSON))
		}

		columns = append(columns, handler.NewCol(UpdatedAt, idpEvent.CreationDate()))

		return handler.NewUpdateStatement(
			&idpEvent,
			columns,
			[]handler.Condition{
				handler.NewCond(IDPTemplateIDCol, idpEvent.ID),
				handler.NewCond(IDPTemplateInstanceIDCol, idpEvent.Aggregate().InstanceID),
				orgCond,
			},
		).Execute(ctx, ex, projectionName)

	}), nil
}

func (p *relationalTablesProjection) reduceIDPChangedTemplateColumns(name *string, optionChanges idp.OptionChanges) []handler.Column {
	cols := make([]handler.Column, 0, 8)
	if name != nil {
		cols = append(cols, handler.NewCol(IDPTemplateNameCol, *name))
	}
	if optionChanges.IsCreationAllowed != nil {
		cols = append(cols, handler.NewCol(IDPAllowCreationCol, *optionChanges.IsCreationAllowed))
	}
	if optionChanges.IsLinkingAllowed != nil {
		cols = append(cols, handler.NewCol(IDPAllowLinkingCol, *optionChanges.IsLinkingAllowed))
	}
	if optionChanges.IsAutoCreation != nil {
		cols = append(cols, handler.NewCol(IDPAllowAutoCreationCol, *optionChanges.IsAutoCreation))
	}
	if optionChanges.IsAutoUpdate != nil {
		cols = append(cols, handler.NewCol(IDPAllowAutoUpdateCol, *optionChanges.IsAutoUpdate))
	}
	if optionChanges.AutoLinkingOption != nil && *optionChanges.AutoLinkingOption != internal_domain.AutoLinkingOptionUnspecified {
		cols = append(cols, handler.NewCol(IDPAllowAutoLinkingCol, domain.IDPAutoLinkingField(*optionChanges.AutoLinkingOption)))
	}

	return cols
}

func (p *relationalTablesProjection) reduceOAuthIDPChangedColumns(payload *domain.OAuth, idpEvent *idp.OAuthIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.AuthorizationEndpoint != nil {
		payloadChange = true
		payload.AuthorizationEndpoint = *idpEvent.AuthorizationEndpoint
	}
	if idpEvent.TokenEndpoint != nil {
		payloadChange = true
		payload.TokenEndpoint = *idpEvent.TokenEndpoint
	}
	if idpEvent.UserEndpoint != nil {
		payloadChange = true
		payload.UserEndpoint = *idpEvent.UserEndpoint
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.IDAttribute != nil {
		payloadChange = true
		payload.IDAttribute = *idpEvent.IDAttribute
	}
	if idpEvent.UsePKCE != nil {
		payloadChange = true
		payload.UsePKCE = *idpEvent.UsePKCE
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceOIDCIDPChangedColumns(payload *domain.OIDC, idpEvent *idp.OIDCIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Issuer != nil {
		payloadChange = true
		payload.Issuer = *idpEvent.Issuer
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.IsIDTokenMapping != nil {
		payloadChange = true
		payload.IsIDTokenMapping = *idpEvent.IsIDTokenMapping
	}
	if idpEvent.UsePKCE != nil {
		payloadChange = true
		payload.UsePKCE = *idpEvent.UsePKCE
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceJWTIDPChangedColumns(payload *domain.JWT, idpEvent *idp.JWTIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.JWTEndpoint != nil {
		payloadChange = true
		payload.JWTEndpoint = *idpEvent.JWTEndpoint
	}
	if idpEvent.KeysEndpoint != nil {
		payloadChange = true
		payload.KeysEndpoint = *idpEvent.KeysEndpoint
	}
	if idpEvent.HeaderName != nil {
		payloadChange = true
		payload.HeaderName = *idpEvent.HeaderName
	}
	if idpEvent.Issuer != nil {
		payloadChange = true
		payload.Issuer = *idpEvent.Issuer
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceAzureADIDPChangedColumns(payload *domain.Azure, idpEvent *idp.AzureADIDPChangedEvent) (bool, error) {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.Tenant != nil {
		payloadChange = true

		azureTenant, err := domain.AzureTenantTypeString(*idpEvent.Tenant)
		if err != nil {
			return false, err
		}

		payload.Tenant = azureTenant
	}
	if idpEvent.IsEmailVerified != nil {
		payloadChange = true
		payload.IsEmailVerified = *idpEvent.IsEmailVerified
	}
	return payloadChange, nil
}

func (p *relationalTablesProjection) reduceGitHubIDPChangedColumns(payload *domain.Github, idpEvent *idp.GitHubIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceGitHubEnterpriseIDPChangedColumns(payload *domain.GithubEnterprise, idpEvent *idp.GitHubEnterpriseIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.AuthorizationEndpoint != nil {
		payloadChange = true
		payload.AuthorizationEndpoint = *idpEvent.AuthorizationEndpoint
	}
	if idpEvent.TokenEndpoint != nil {
		payloadChange = true
		payload.TokenEndpoint = *idpEvent.TokenEndpoint
	}
	if idpEvent.UserEndpoint != nil {
		payloadChange = true
		payload.UserEndpoint = *idpEvent.UserEndpoint
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceGitLabIDPChangedColumns(payload *domain.Gitlab, idpEvent *idp.GitLabIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceGitLabSelfHostedIDPChangedColumns(payload *domain.GitlabSelfHosting, idpEvent *idp.GitLabSelfHostedIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Issuer != nil {
		payloadChange = true
		payload.Issuer = *idpEvent.Issuer
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceGoogleIDPChangedColumns(payload *domain.Google, idpEvent *idp.GoogleIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChange = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceLDAPIDPChangedColumns(payload *domain.LDAP, idpEvent *idp.LDAPIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.Servers != nil {
		payloadChange = true
		payload.Servers = idpEvent.Servers
	}
	if idpEvent.StartTLS != nil {
		payloadChange = true
		payload.StartTLS = *idpEvent.StartTLS
	}
	if idpEvent.BaseDN != nil {
		payloadChange = true
		payload.BaseDN = *idpEvent.BaseDN
	}
	if idpEvent.BindDN != nil {
		payloadChange = true
		payload.BindDN = *idpEvent.BindDN
	}
	if idpEvent.BindPassword != nil {
		payloadChange = true
		payload.BindPassword = idpEvent.BindPassword
	}
	if idpEvent.UserBase != nil {
		payloadChange = true
		payload.UserBase = *idpEvent.UserBase
	}
	if idpEvent.UserObjectClasses != nil {
		payloadChange = true
		payload.UserObjectClasses = idpEvent.UserObjectClasses
	}
	if idpEvent.UserFilters != nil {
		payloadChange = true
		payload.UserFilters = idpEvent.UserFilters
	}
	if idpEvent.Timeout != nil {
		payloadChange = true
		payload.Timeout = *idpEvent.Timeout
	}
	if idpEvent.RootCA != nil {
		payloadChange = true
		payload.RootCA = idpEvent.RootCA
	}
	if idpEvent.IDAttribute != nil {
		payloadChange = true
		payload.IDAttribute = *idpEvent.IDAttribute
	}
	if idpEvent.FirstNameAttribute != nil {
		payloadChange = true
		payload.FirstNameAttribute = *idpEvent.FirstNameAttribute
	}
	if idpEvent.LastNameAttribute != nil {
		payloadChange = true
		payload.LastNameAttribute = *idpEvent.LastNameAttribute
	}
	if idpEvent.DisplayNameAttribute != nil {
		payloadChange = true
		payload.DisplayNameAttribute = *idpEvent.DisplayNameAttribute
	}
	if idpEvent.NickNameAttribute != nil {
		payloadChange = true
		payload.NickNameAttribute = *idpEvent.NickNameAttribute
	}
	if idpEvent.PreferredUsernameAttribute != nil {
		payloadChange = true
		payload.PreferredUsernameAttribute = *idpEvent.PreferredUsernameAttribute
	}
	if idpEvent.EmailAttribute != nil {
		payloadChange = true
		payload.EmailAttribute = *idpEvent.EmailAttribute
	}
	if idpEvent.EmailVerifiedAttribute != nil {
		payloadChange = true
		payload.EmailVerifiedAttribute = *idpEvent.EmailVerifiedAttribute
	}
	if idpEvent.PhoneAttribute != nil {
		payloadChange = true
		payload.PhoneAttribute = *idpEvent.PhoneAttribute
	}
	if idpEvent.PhoneVerifiedAttribute != nil {
		payloadChange = true
		payload.PhoneVerifiedAttribute = *idpEvent.PhoneVerifiedAttribute
	}
	if idpEvent.PreferredLanguageAttribute != nil {
		payloadChange = true
		payload.PreferredLanguageAttribute = *idpEvent.PreferredLanguageAttribute
	}
	if idpEvent.AvatarURLAttribute != nil {
		payloadChange = true
		payload.AvatarURLAttribute = *idpEvent.AvatarURLAttribute
	}
	if idpEvent.ProfileAttribute != nil {
		payloadChange = true
		payload.ProfileAttribute = *idpEvent.ProfileAttribute
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceAppleIDPChangedColumns(payload *domain.Apple, idpEvent *idp.AppleIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.ClientID != nil {
		payloadChange = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.TeamID != nil {
		payloadChange = true
		payload.TeamID = *idpEvent.TeamID
	}
	if idpEvent.KeyID != nil {
		payloadChange = true
		payload.KeyID = *idpEvent.KeyID
	}
	if idpEvent.PrivateKey != nil {
		payloadChange = true
		payload.PrivateKey = idpEvent.PrivateKey
	}
	if idpEvent.Scopes != nil {
		payloadChange = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChange
}

func (p *relationalTablesProjection) reduceSAMLIDPChangedColumns(payload *domain.SAML, idpEvent *idp.SAMLIDPChangedEvent) bool {
	payloadChange := false
	if idpEvent.Metadata != nil {
		payloadChange = true
		payload.Metadata = idpEvent.Metadata
	}
	if idpEvent.Key != nil {
		payloadChange = true
		payload.Key = idpEvent.Key
	}
	if idpEvent.Certificate != nil {
		payloadChange = true
		payload.Certificate = idpEvent.Certificate
	}
	if idpEvent.Binding != nil {
		payloadChange = true
		payload.Binding = *idpEvent.Binding
	}
	if idpEvent.WithSignedRequest != nil {
		payloadChange = true
		payload.WithSignedRequest = *idpEvent.WithSignedRequest
	}
	if idpEvent.NameIDFormat != nil {
		payloadChange = true
		payload.NameIDFormat = idpEvent.NameIDFormat
	}
	if idpEvent.TransientMappingAttributeName != nil {
		payloadChange = true
		payload.TransientMappingAttributeName = *idpEvent.TransientMappingAttributeName
	}
	if idpEvent.FederatedLogoutEnabled != nil {
		payloadChange = true
		payload.FederatedLogoutEnabled = *idpEvent.FederatedLogoutEnabled
	}
	if idpEvent.SignatureAlgorithm != nil {
		payloadChange = true
		payload.SignatureAlgorithm = *idpEvent.SignatureAlgorithm
	}
	return payloadChange
}
