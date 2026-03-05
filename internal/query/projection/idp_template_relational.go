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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        event.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ConfigID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              mapIDPConfigType(idpEvent.Typ),
			AllowAutoCreation: idpEvent.AutoRegister,
			AllowLinking:      true,
			AllowCreation:     true,
			CreatedAt:         event.CreatedAt(),
			UpdatedAt:         event.CreatedAt(),
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
	case internal_domain.IDPConfigTypeUnspecified:
		fallthrough
	default:
		return nil
	}
}

func mapAutoLinkingField(option internal_domain.AutoLinkingOption) *domain.IDPAutoLinkingField {
	if option == internal_domain.AutoLinkingOptionUnspecified {
		return nil
	}
	switch option {
	case internal_domain.AutoLinkingOptionEmail:
		return gu.Ptr(domain.IDPAutoLinkingFieldEmail)
	case internal_domain.AutoLinkingOptionUsername:
		return gu.Ptr(domain.IDPAutoLinkingFieldUsername)
	case internal_domain.AutoLinkingOptionUnspecified:
		fallthrough
	default:
		return nil
	}
}

func idpScopedCondition(repo domain.IDProviderRepository, instanceID, id string, orgID *string) database.Condition {
	return database.And(
		repo.PrimaryKeyCondition(instanceID, id),
		repo.OrgIDCondition(orgID),
	)
}

func (p *relationalTablesProjection) reduceIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgID *string
	var idpEvent idpconfig.IDPConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
		orgID = gu.Ptr(idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPConfigChangedEvent:
		idpEvent = e.IDPConfigChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YVvJD", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigChangedEventType, instance.IDPConfigChangedEventType})
	}

	repo := repository.IDProviderRepository()

	changes := make(database.Changes, 0, 3)
	if idpEvent.Name != nil {
		changes = append(changes, repo.SetName(*idpEvent.Name))
	}
	if idpEvent.AutoRegister != nil {
		changes = append(changes, repo.SetAllowAutoCreation(*idpEvent.AutoRegister))
	}
	if len(changes) == 0 {
		return handler.NewNoOpStatement(&idpEvent), nil
	}
	changes = append(changes, repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())))

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-9sX8h", "reduce.wrong.db.pool %T", ex)
		}

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, event.Aggregate().InstanceID, idpEvent.ConfigID, orgID), changes...)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceIDPDeactivated(event eventstore.Event) (*handler.Statement, error) {
	var orgID *string
	var idpEvent idpconfig.IDPConfigDeactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
		orgID = gu.Ptr(idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPConfigDeactivatedEvent:
		idpEvent = e.IDPConfigDeactivatedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y4O5l", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigDeactivatedEventType, instance.IDPConfigDeactivatedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-9sX8h", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, event.Aggregate().InstanceID, idpEvent.ConfigID, orgID),
			repo.SetState(domain.IDPStateInactive),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceIDPReactivated(event eventstore.Event) (*handler.Statement, error) {
	var orgID *string
	var idpEvent idpconfig.IDPConfigReactivatedEvent
	switch e := event.(type) {
	case *org.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
		orgID = gu.Ptr(idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPConfigReactivatedEvent:
		idpEvent = e.IDPConfigReactivatedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y8QyS", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPConfigReactivatedEventType, instance.IDPConfigReactivatedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-2Db9P", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.ConfigID, orgID),
			repo.SetState(domain.IDPStateActive),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
	var (
		orgID *string
		idpID string
	)
	switch e := event.(type) {
	case *org.IDPRemovedEvent:
		idpID = e.ID
		orgID = gu.Ptr(e.RemovedEvent.Aggregate().ResourceOwner)
	case *instance.IDPRemovedEvent:
		idpID = e.ID
	case *org.IDPConfigRemovedEvent:
		idpID = e.ConfigID
		orgID = gu.Ptr(e.IDPConfigRemovedEvent.Aggregate().ResourceOwner)
	case *instance.IDPConfigRemovedEvent:
		idpID = e.ConfigID
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ybcvwin2", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPRemovedEventType, instance.IDPRemovedEventType, org.IDPConfigRemovedEventType, instance.IDPConfigRemovedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-PSj7F", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, event.Aggregate().InstanceID, idpID, orgID))
		return err
	}), nil
}

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

func (p *relationalTablesProjection) reduceJWTConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	var orgID *string
	var idpEvent idpconfig.JWTConfigAddedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
		orgID = gu.Ptr(idpEvent.Aggregate().ResourceOwner)
	case *instance.IDPJWTConfigAddedEvent:
		idpEvent = e.JWTConfigAddedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YvPdb", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigAddedEventType, instance.IDPJWTConfigAddedEventType})
	}

	payloadJSON, err := json.Marshal(idpEvent)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-tJQ8V", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		_, err = repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.IDPConfigID, orgID),
			repo.SetPayload(string(payloadJSON)),
			repo.SetType(domain.IDPTypeJWT),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idpconfig.JWTConfigChangedEvent
	switch e := event.(type) {
	case *org.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.IDPJWTConfigChangedEvent:
		idpEvent = e.JWTConfigChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y2IVI", "reduce.wrong.event.type %v", []eventstore.EventType{org.IDPJWTConfigChangedEventType, instance.IDPJWTConfigChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-sh6Lp", "unable to cast to tx executer")
		}

		idpRepo := repository.IDProviderRepository()
		jwt, err := idpRepo.GetJWT(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.IDPConfigID, orgId)))
		if err != nil {
			return err
		}

		if idpEvent.JWTEndpoint != nil && *idpEvent.JWTEndpoint != jwt.JWTEndpoint {
			jwt.JWTEndpoint = *idpEvent.JWTEndpoint
		}
		if idpEvent.Issuer != nil && *idpEvent.Issuer != jwt.Issuer {
			jwt.Issuer = *idpEvent.Issuer
		}
		if idpEvent.KeysEndpoint != nil && *idpEvent.KeysEndpoint != jwt.KeysEndpoint {
			jwt.KeysEndpoint = *idpEvent.KeysEndpoint
		}
		if idpEvent.HeaderName != nil && *idpEvent.HeaderName != jwt.HeaderName {
			jwt.HeaderName = *idpEvent.HeaderName
		}

		payloadJSON, err := json.Marshal(jwt.JWT)
		if err != nil {
			return err
		}

		_, err = idpRepo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(idpRepo, idpEvent.Aggregate().InstanceID, idpEvent.IDPConfigID, orgId),
			idpRepo.SetPayload(string(payloadJSON)),
			idpRepo.SetType(domain.IDPTypeJWT),
			idpRepo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
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

func (p *relationalTablesProjection) reduceOIDCIDPMigratedAzureAD(event eventstore.Event) (*handler.Statement, error) {
	var orgID *string
	var idpEvent idp.OIDCIDPMigratedAzureADEvent
	switch e := event.(type) {
	case *org.OIDCIDPMigratedAzureADEvent:
		idpEvent = e.OIDCIDPMigratedAzureADEvent
		orgID = gu.Ptr(idpEvent.Aggregate().ResourceOwner)
	case *instance.OIDCIDPMigratedAzureADEvent:
		idpEvent = e.OIDCIDPMigratedAzureADEvent
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-mj7LQ", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		changes := database.Changes{
			repo.SetName(idpEvent.Name),
			repo.SetType(domain.IDPTypeAzure),
			repo.SetAllowCreation(idpEvent.IsCreationAllowed),
			repo.SetAllowLinking(idpEvent.IsLinkingAllowed),
			repo.SetAllowAutoCreation(idpEvent.IsAutoCreation),
			repo.SetAllowAutoUpdate(idpEvent.IsAutoUpdate),
			repo.SetAutoLinkingField(mapAutoLinkingField(idpEvent.AutoLinkingOption)),
			repo.SetPayload(string(payloadJSON)),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		}

		_, err = repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.ID, orgID), changes...)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceOIDCIDPMigratedGoogle(event eventstore.Event) (*handler.Statement, error) {
	var orgID *string
	var idpEvent idp.OIDCIDPMigratedGoogleEvent
	switch e := event.(type) {
	case *org.OIDCIDPMigratedGoogleEvent:
		idpEvent = e.OIDCIDPMigratedGoogleEvent
		orgID = gu.Ptr(idpEvent.Aggregate().ResourceOwner)
	case *instance.OIDCIDPMigratedGoogleEvent:
		idpEvent = e.OIDCIDPMigratedGoogleEvent
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-HDqk9", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		changes := database.Changes{
			repo.SetName(idpEvent.Name),
			repo.SetType(domain.IDPTypeGoogle),
			repo.SetAllowCreation(idpEvent.IsCreationAllowed),
			repo.SetAllowLinking(idpEvent.IsLinkingAllowed),
			repo.SetAllowAutoCreation(idpEvent.IsAutoCreation),
			repo.SetAllowAutoUpdate(idpEvent.IsAutoUpdate),
			repo.SetAutoLinkingField(mapAutoLinkingField(idpEvent.AutoLinkingOption)),
			repo.SetPayload(string(payloadJSON)),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		}

		_, err = repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.ID, orgID), changes...)
		return err
	}), nil
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-ZYYyQ", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeJWT),
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

func (p *relationalTablesProjection) reduceJWTIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.JWTIDPChangedEvent
	switch e := event.(type) {
	case *org.JWTIDPChangedEvent:
		idpEvent = e.JWTIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.JWTIDPChangedEvent:
		idpEvent = e.JWTIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-H15j2il", "reduce.wrong.event.type %v", []eventstore.EventType{org.JWTIDPChangedEventType, instance.JWTIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		jwt, err := idpRepo.GetJWT(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &jwt.JWT
		payloadChanged := p.reduceJWTIDPChangedColumns(payload, &idpEvent)
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-GJ4Kb", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeAzure),
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

func (p *relationalTablesProjection) reduceAzureADIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.AzureADIDPChangedEvent
	switch e := event.(type) {
	case *org.AzureADIDPChangedEvent:
		idpEvent = e.AzureADIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.AzureADIDPChangedEvent:
		idpEvent = e.AzureADIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YZ5x25s", "reduce.wrong.event.type %v", []eventstore.EventType{org.AzureADIDPChangedEventType, instance.AzureADIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		azure, err := idpRepo.GetAzureAD(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

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
			changes = append(changes, idpRepo.SetPayload(string(payloadJSON)))
		}

		changes = append(changes, idpRepo.SetUpdatedAt(gu.Ptr(event.CreatedAt())))
		_, err = idpRepo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(idpRepo, idpEvent.Aggregate().InstanceID, idpEvent.ID, orgId), changes...)
		return err

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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-HNpgd", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeGitHub),
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

func (p *relationalTablesProjection) reduceGitHubIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GitHubIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GitHubIDPChangedEvent:
		idpEvent = e.GitHubIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-L1U89ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubIDPChangedEventType, instance.GitHubIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		github, err := idpRepo.GetGithub(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &github.Github
		payloadChanged := p.reduceGitHubIDPChangedColumns(payload, &idpEvent)
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-Kv4Fu", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeGitHubEnterprise),
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

func (p *relationalTablesProjection) reduceGitHubEnterpriseIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GitHubEnterpriseIDPChangedEvent
	switch e := event.(type) {
	case *org.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GitHubEnterpriseIDPChangedEvent:
		idpEvent = e.GitHubEnterpriseIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YDg3g", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitHubEnterpriseIDPChangedEventType, instance.GitHubEnterpriseIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		githubEnterprise, err := idpRepo.GetGithubEnterprise(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &githubEnterprise.GithubEnterprise
		payloadChanged := p.reduceGitHubEnterpriseIDPChangedColumns(payload, &idpEvent)
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kN8Qx", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeGitLab),
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

func (p *relationalTablesProjection) reduceGitLabIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GitLabIDPChangedEvent
	switch e := event.(type) {
	case *org.GitLabIDPChangedEvent:
		idpEvent = e.GitLabIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GitLabIDPChangedEvent:
		idpEvent = e.GitLabIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-mT5827b", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabIDPChangedEventType, instance.GitLabIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		gitlab, err := idpRepo.GetGitlab(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &gitlab.Gitlab
		payloadChanged := p.reduceGitLabIDPChangedColumns(payload, &idpEvent)
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

	gitlabSelfHosted := domain.GitlabSelfHosted{
		Issuer:       idpEvent.Issuer,
		ClientID:     idpEvent.ClientID,
		ClientSecret: idpEvent.ClientSecret,
		Scopes:       idpEvent.Scopes,
	}

	payloadJSON, err := json.Marshal(gitlabSelfHosted)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FQrtw", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeGitLabSelfHosted),
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

func (p *relationalTablesProjection) reduceGitLabSelfHostedIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GitLabSelfHostedIDPChangedEvent
	switch e := event.(type) {
	case *org.GitLabSelfHostedIDPChangedEvent:
		idpEvent = e.GitLabSelfHostedIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GitLabSelfHostedIDPChangedEvent:
		idpEvent = e.GitLabSelfHostedIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YAf3g2", "reduce.wrong.event.type %v", []eventstore.EventType{org.GitLabSelfHostedIDPChangedEventType, instance.GitLabSelfHostedIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		gitlabSelfHosted, err := idpRepo.GetGitlabSelfHosted(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &gitlabSelfHosted.GitlabSelfHosted
		payloadChanged := p.reduceGitLabSelfHostedIDPChangedColumns(payload, &idpEvent)
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-B9SPm", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeGoogle),
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

func (p *relationalTablesProjection) reduceGoogleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.GoogleIDPChangedEvent
	switch e := event.(type) {
	case *org.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.GoogleIDPChangedEvent:
		idpEvent = e.GoogleIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YN58hml", "reduce.wrong.event.type %v", []eventstore.EventType{org.GoogleIDPChangedEventType, instance.GoogleIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		google, err := idpRepo.GetGoogle(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &google.Google
		payloadChanged := p.reduceGoogleIDPChangedColumns(payload, &idpEvent)
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-XCJ8w", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeLDAP),
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

func (p *relationalTablesProjection) reduceLDAPIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.LDAPIDPChangedEvent
	switch e := event.(type) {
	case *org.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.LDAPIDPChangedEvent:
		idpEvent = e.LDAPIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-p1582ks", "reduce.wrong.event.type %v", []eventstore.EventType{org.LDAPIDPChangedEventType, instance.LDAPIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		ldap, err := idpRepo.GetLDAP(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &ldap.LDAP
		payloadChanged := p.reduceLDAPIDPChangedColumns(payload, &idpEvent)
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ku2YB", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeApple),
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

func (p *relationalTablesProjection) reduceAppleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.AppleIDPChangedEvent
	switch e := event.(type) {
	case *org.AppleIDPChangedEvent:
		idpEvent = e.AppleIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.AppleIDPChangedEvent:
		idpEvent = e.AppleIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YBez3", "reduce.wrong.event.type %v", []eventstore.EventType{org.AppleIDPChangedEventType /*, instance.AppleIDPChangedEventType*/})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		apple, err := idpRepo.GetApple(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &apple.Apple
		payloadChanged := p.reduceAppleIDPChangedColumns(payload, &idpEvent)
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-ksJ3N", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &domain.IdentityProvider{
			InstanceID:        idpEvent.Aggregate().InstanceID,
			OrgID:             orgId,
			ID:                idpEvent.ID,
			State:             domain.IDPStateActive,
			Name:              idpEvent.Name,
			Type:              gu.Ptr(domain.IDPTypeSAML),
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

func (p *relationalTablesProjection) reduceSAMLIDPChanged(event eventstore.Event) (*handler.Statement, error) {
	var orgId *string
	var idpEvent idp.SAMLIDPChangedEvent
	switch e := event.(type) {
	case *org.SAMLIDPChangedEvent:
		idpEvent = e.SAMLIDPChangedEvent
		orgId = &idpEvent.Aggregate().ResourceOwner
	case *instance.SAMLIDPChangedEvent:
		idpEvent = e.SAMLIDPChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Y7c0fii4ad", "reduce.wrong.event.type %v", []eventstore.EventType{org.SAMLIDPChangedEventType, instance.SAMLIDPChangedEventType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-HX6ed", "unable to cast to tx executer")
		}
		idpRepo := repository.IDProviderRepository()
		saml, err := idpRepo.GetSAML(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.ID, orgId)))
		if err != nil {
			return err
		}

		changes := p.reduceIDPChangedTemplateColumns(idpRepo, idpEvent.Name, idpEvent.OptionChanges)

		payload := &saml.SAML
		payloadChanged := p.reduceSAMLIDPChangedColumns(payload, &idpEvent)
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

func (p *relationalTablesProjection) reduceIDPChangedTemplateColumns(repo domain.IDProviderRepository, name *string, optionChanges idp.OptionChanges) database.Changes {
	changes := make(database.Changes, 0, 8)
	if name != nil {
		changes = append(changes, repo.SetName(*name))
	}
	if optionChanges.IsCreationAllowed != nil {
		changes = append(changes, repo.SetAllowCreation(*optionChanges.IsCreationAllowed))
	}
	if optionChanges.IsLinkingAllowed != nil {
		changes = append(changes, repo.SetAllowLinking(*optionChanges.IsLinkingAllowed))
	}
	if optionChanges.IsAutoCreation != nil {
		changes = append(changes, repo.SetAllowAutoCreation(*optionChanges.IsAutoCreation))
	}
	if optionChanges.IsAutoUpdate != nil {
		changes = append(changes, repo.SetAllowAutoUpdate(*optionChanges.IsAutoUpdate))
	}
	if optionChanges.AutoLinkingOption != nil {
		changes = append(changes, repo.SetAutoLinkingField(mapAutoLinkingField(*optionChanges.AutoLinkingOption)))
	}

	return changes
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

func (p *relationalTablesProjection) reduceJWTIDPChangedColumns(payload *domain.JWT, idpEvent *idp.JWTIDPChangedEvent) bool {
	payloadChanged := false
	if payload == nil || idpEvent == nil {
		return payloadChanged
	}

	if idpEvent.JWTEndpoint != nil && *idpEvent.JWTEndpoint != payload.JWTEndpoint {
		payloadChanged = true
		payload.JWTEndpoint = *idpEvent.JWTEndpoint
	}
	if idpEvent.KeysEndpoint != nil && *idpEvent.KeysEndpoint != payload.KeysEndpoint {
		payloadChanged = true
		payload.KeysEndpoint = *idpEvent.KeysEndpoint
	}
	if idpEvent.HeaderName != nil && *idpEvent.HeaderName != payload.HeaderName {
		payloadChanged = true
		payload.HeaderName = *idpEvent.HeaderName
	}
	if idpEvent.Issuer != nil && *idpEvent.Issuer != payload.Issuer {
		payloadChanged = true
		payload.Issuer = *idpEvent.Issuer
	}
	return payloadChanged
}

func (p *relationalTablesProjection) reduceAzureADIDPChangedColumns(payload *domain.Azure, idpEvent *idp.AzureADIDPChangedEvent) (bool, error) {
	payloadChanged := false
	if payload == nil || idpEvent == nil {
		return payloadChanged, nil
	}

	if idpEvent.ClientID != nil && *idpEvent.ClientID != payload.ClientID {
		payloadChanged = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.ClientSecret != nil {
		payloadChanged = true
		payload.ClientSecret = idpEvent.ClientSecret
	}
	if idpEvent.Scopes != nil && !slices.Equal(idpEvent.Scopes, payload.Scopes) {
		payloadChanged = true
		payload.Scopes = idpEvent.Scopes
	}
	if idpEvent.Tenant != nil && *idpEvent.Tenant != payload.Tenant.String() {
		payloadChanged = true

		azureTenant, err := domain.AzureTenantTypeString(*idpEvent.Tenant)
		if err != nil {
			return false, err
		}

		payload.Tenant = azureTenant
	}
	if idpEvent.IsEmailVerified != nil && *idpEvent.IsEmailVerified != payload.IsEmailVerified {
		payloadChanged = true
		payload.IsEmailVerified = *idpEvent.IsEmailVerified
	}
	return payloadChanged, nil
}

func (p *relationalTablesProjection) reduceGitHubIDPChangedColumns(payload *domain.Github, idpEvent *idp.GitHubIDPChangedEvent) bool {
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
	if idpEvent.Scopes != nil && !slices.Equal(idpEvent.Scopes, payload.Scopes) {
		payloadChanged = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChanged
}

func (p *relationalTablesProjection) reduceGitHubEnterpriseIDPChangedColumns(payload *domain.GithubEnterprise, idpEvent *idp.GitHubEnterpriseIDPChangedEvent) bool {
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
	return payloadChanged
}

func (p *relationalTablesProjection) reduceGitLabIDPChangedColumns(payload *domain.Gitlab, idpEvent *idp.GitLabIDPChangedEvent) bool {
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
	if idpEvent.Scopes != nil && !slices.Equal(idpEvent.Scopes, payload.Scopes) {
		payloadChanged = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChanged
}

func (p *relationalTablesProjection) reduceGitLabSelfHostedIDPChangedColumns(payload *domain.GitlabSelfHosted, idpEvent *idp.GitLabSelfHostedIDPChangedEvent) bool {
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
	return payloadChanged
}

func (p *relationalTablesProjection) reduceGoogleIDPChangedColumns(payload *domain.Google, idpEvent *idp.GoogleIDPChangedEvent) bool {
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
	if idpEvent.Scopes != nil && !slices.Equal(idpEvent.Scopes, payload.Scopes) {
		payloadChanged = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChanged
}

//nolint:gocognit
func (p *relationalTablesProjection) reduceLDAPIDPChangedColumns(payload *domain.LDAP, idpEvent *idp.LDAPIDPChangedEvent) bool {
	payloadChanged := false
	if payload == nil || idpEvent == nil {
		return payloadChanged
	}

	if idpEvent.Servers != nil && !slices.Equal(idpEvent.Servers, payload.Servers) {
		payloadChanged = true
		payload.Servers = idpEvent.Servers
	}
	if idpEvent.StartTLS != nil && *idpEvent.StartTLS != payload.StartTLS {
		payloadChanged = true
		payload.StartTLS = *idpEvent.StartTLS
	}
	if idpEvent.BaseDN != nil && *idpEvent.BaseDN != payload.BaseDN {
		payloadChanged = true
		payload.BaseDN = *idpEvent.BaseDN
	}
	if idpEvent.BindDN != nil && *idpEvent.BindDN != payload.BindDN {
		payloadChanged = true
		payload.BindDN = *idpEvent.BindDN
	}
	if idpEvent.BindPassword != nil {
		payloadChanged = true
		payload.BindPassword = idpEvent.BindPassword
	}
	if idpEvent.UserBase != nil && *idpEvent.UserBase != payload.UserBase {
		payloadChanged = true
		payload.UserBase = *idpEvent.UserBase
	}
	if idpEvent.UserObjectClasses != nil && !slices.Equal(idpEvent.UserObjectClasses, payload.UserObjectClasses) {
		payloadChanged = true
		payload.UserObjectClasses = idpEvent.UserObjectClasses
	}
	if idpEvent.UserFilters != nil && !slices.Equal(idpEvent.UserFilters, payload.UserFilters) {
		payloadChanged = true
		payload.UserFilters = idpEvent.UserFilters
	}
	if idpEvent.Timeout != nil && *idpEvent.Timeout != payload.Timeout {
		payloadChanged = true
		payload.Timeout = *idpEvent.Timeout
	}
	if idpEvent.RootCA != nil && !slices.Equal(idpEvent.RootCA, payload.RootCA) {
		payloadChanged = true
		payload.RootCA = idpEvent.RootCA
	}
	if idpEvent.IDAttribute != nil && *idpEvent.IDAttribute != payload.IDAttribute {
		payloadChanged = true
		payload.IDAttribute = *idpEvent.IDAttribute
	}
	if idpEvent.FirstNameAttribute != nil && *idpEvent.FirstNameAttribute != payload.FirstNameAttribute {
		payloadChanged = true
		payload.FirstNameAttribute = *idpEvent.FirstNameAttribute
	}
	if idpEvent.LastNameAttribute != nil && *idpEvent.LastNameAttribute != payload.LastNameAttribute {
		payloadChanged = true
		payload.LastNameAttribute = *idpEvent.LastNameAttribute
	}
	if idpEvent.DisplayNameAttribute != nil && *idpEvent.DisplayNameAttribute != payload.DisplayNameAttribute {
		payloadChanged = true
		payload.DisplayNameAttribute = *idpEvent.DisplayNameAttribute
	}
	if idpEvent.NickNameAttribute != nil && *idpEvent.NickNameAttribute != payload.NickNameAttribute {
		payloadChanged = true
		payload.NickNameAttribute = *idpEvent.NickNameAttribute
	}
	if idpEvent.PreferredUsernameAttribute != nil && *idpEvent.PreferredUsernameAttribute != payload.PreferredUsernameAttribute {
		payloadChanged = true
		payload.PreferredUsernameAttribute = *idpEvent.PreferredUsernameAttribute
	}
	if idpEvent.EmailAttribute != nil && *idpEvent.EmailAttribute != payload.EmailAttribute {
		payloadChanged = true
		payload.EmailAttribute = *idpEvent.EmailAttribute
	}
	if idpEvent.EmailVerifiedAttribute != nil && *idpEvent.EmailVerifiedAttribute != payload.EmailVerifiedAttribute {
		payloadChanged = true
		payload.EmailVerifiedAttribute = *idpEvent.EmailVerifiedAttribute
	}
	if idpEvent.PhoneAttribute != nil && *idpEvent.PhoneAttribute != payload.PhoneAttribute {
		payloadChanged = true
		payload.PhoneAttribute = *idpEvent.PhoneAttribute
	}
	if idpEvent.PhoneVerifiedAttribute != nil && *idpEvent.PhoneVerifiedAttribute != payload.PhoneVerifiedAttribute {
		payloadChanged = true
		payload.PhoneVerifiedAttribute = *idpEvent.PhoneVerifiedAttribute
	}
	if idpEvent.PreferredLanguageAttribute != nil && *idpEvent.PreferredLanguageAttribute != payload.PreferredLanguageAttribute {
		payloadChanged = true
		payload.PreferredLanguageAttribute = *idpEvent.PreferredLanguageAttribute
	}
	if idpEvent.AvatarURLAttribute != nil && *idpEvent.AvatarURLAttribute != payload.AvatarURLAttribute {
		payloadChanged = true
		payload.AvatarURLAttribute = *idpEvent.AvatarURLAttribute
	}
	if idpEvent.ProfileAttribute != nil && *idpEvent.ProfileAttribute != payload.ProfileAttribute {
		payloadChanged = true
		payload.ProfileAttribute = *idpEvent.ProfileAttribute
	}
	return payloadChanged
}

func (p *relationalTablesProjection) reduceAppleIDPChangedColumns(payload *domain.Apple, idpEvent *idp.AppleIDPChangedEvent) bool {
	payloadChanged := false
	if payload == nil || idpEvent == nil {
		return payloadChanged
	}

	if idpEvent.ClientID != nil && *idpEvent.ClientID != payload.ClientID {
		payloadChanged = true
		payload.ClientID = *idpEvent.ClientID
	}
	if idpEvent.TeamID != nil && *idpEvent.TeamID != payload.TeamID {
		payloadChanged = true
		payload.TeamID = *idpEvent.TeamID
	}
	if idpEvent.KeyID != nil && *idpEvent.KeyID != payload.KeyID {
		payloadChanged = true
		payload.KeyID = *idpEvent.KeyID
	}
	if idpEvent.PrivateKey != nil {
		payloadChanged = true
		payload.PrivateKey = idpEvent.PrivateKey
	}
	if idpEvent.Scopes != nil && !slices.Equal(idpEvent.Scopes, payload.Scopes) {
		payloadChanged = true
		payload.Scopes = idpEvent.Scopes
	}
	return payloadChanged
}

func (p *relationalTablesProjection) reduceSAMLIDPChangedColumns(payload *domain.SAML, idpEvent *idp.SAMLIDPChangedEvent) bool {
	payloadChanged := false
	if payload == nil || idpEvent == nil {
		return payloadChanged
	}

	if idpEvent.Metadata != nil && !slices.Equal(idpEvent.Metadata, payload.Metadata) {
		payloadChanged = true
		payload.Metadata = idpEvent.Metadata
	}
	if idpEvent.Key != nil {
		payloadChanged = true
		payload.Key = idpEvent.Key
	}
	if idpEvent.Certificate != nil && !slices.Equal(idpEvent.Certificate, payload.Certificate) {
		payloadChanged = true
		payload.Certificate = idpEvent.Certificate
	}
	if idpEvent.Binding != nil && *idpEvent.Binding != payload.Binding {
		payloadChanged = true
		payload.Binding = *idpEvent.Binding
	}
	if idpEvent.WithSignedRequest != nil && *idpEvent.WithSignedRequest != payload.WithSignedRequest {
		payloadChanged = true
		payload.WithSignedRequest = *idpEvent.WithSignedRequest
	}
	if idpEvent.NameIDFormat != nil && (payload.NameIDFormat == nil || *idpEvent.NameIDFormat != *payload.NameIDFormat) {
		payloadChanged = true
		payload.NameIDFormat = idpEvent.NameIDFormat
	}
	if idpEvent.TransientMappingAttributeName != nil && *idpEvent.TransientMappingAttributeName != payload.TransientMappingAttributeName {
		payloadChanged = true
		payload.TransientMappingAttributeName = *idpEvent.TransientMappingAttributeName
	}
	if idpEvent.FederatedLogoutEnabled != nil && *idpEvent.FederatedLogoutEnabled != payload.FederatedLogoutEnabled {
		payloadChanged = true
		payload.FederatedLogoutEnabled = *idpEvent.FederatedLogoutEnabled
	}
	if idpEvent.SignatureAlgorithm != nil && *idpEvent.SignatureAlgorithm != payload.SignatureAlgorithm {
		payloadChanged = true
		payload.SignatureAlgorithm = *idpEvent.SignatureAlgorithm
	}
	return payloadChanged
}
