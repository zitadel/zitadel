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
	IDPRelationalTable                = "zitadel.identity_providers"
	IDPRelationalOrgIdCol             = "org_id"
	IDPRelationalAutoRegisterCol      = "auto_register"
	IDPRelationalPayloadCol           = "payload"
	IDPRelationalOrgId                = "org_id"
	IDPRelationalAllowCreationCol     = "allow_creation"
	IDPRelationalAllowLinkingCol      = "allow_linking"
	IDPRelationalAllowAutoCreationCol = "allow_auto_creation"
	IDPRelationalAllowAutoUpdateCol   = "allow_auto_update"
	IDPRelationalAllowAutoLinkingCol  = "auto_linking_field"
)

type idpTemplateRelationalProjection struct {
	idpRepo domain.IDProviderRepository
}

func newIDPTemplateRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	idpRepo := repository.IDProviderRepository()
	return handler.NewHandler(ctx, &config, &idpTemplateRelationalProjection{
		idpRepo: idpRepo,
	})
}

func (*idpTemplateRelationalProjection) Name() string {
	return IDPRelationalTable
}

func (p *idpTemplateRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.IDPConfigAddedEventType,
					Reduce: p.reduceIDPRelationalAdded,
				},
				{
					Event:  instance.IDPConfigChangedEventType,
					Reduce: p.reduceIDPRelationalChanged,
				},
				{
					Event:  instance.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDRelationalPDeactivated,
				},
				{
					Event:  instance.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPRelationalReactivated,
				},
				{
					Event:  instance.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRelationalRemoved,
				},
				{
					Event:  instance.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCRelationalConfigAdded,
				},
				{
					Event:  instance.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCRelationalConfigChanged,
				},
				{
					Event:  instance.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTRelationalConfigAdded,
				},
				{
					Event:  instance.IDPJWTConfigChangedEventType,
					Reduce: p.reduceJWTRelationalConfigChanged,
				},
				{
					Event:  instance.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPRelationalAdded,
				},
				{
					Event:  instance.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPRelationalChanged,
				},
				{
					Event:  instance.OIDCIDPAddedEventType,
					Reduce: p.reduceOIDCIDPRelationalAdded,
				},
				{
					Event:  instance.OIDCIDPChangedEventType,
					Reduce: p.reduceOIDCIDPRelationalChanged,
				},
				{
					Event:  instance.OIDCIDPMigratedAzureADEventType,
					Reduce: p.reduceOIDCIDPRelationalMigratedAzureAD,
				},
				{
					Event:  instance.OIDCIDPMigratedGoogleEventType,
					Reduce: p.reduceOIDCIDPRelationalMigratedGoogle,
				},
				{
					Event:  instance.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPRelationalAdded,
				},
				{
					Event:  instance.JWTIDPChangedEventType,
					Reduce: p.reduceJWTIDPRelationalChanged,
				},
				{
					Event:  instance.AzureADIDPAddedEventType,
					Reduce: p.reduceAzureADIDPRelationalAdded,
				},
				{
					Event:  instance.AzureADIDPChangedEventType,
					Reduce: p.reduceAzureADIDPRelationalChanged,
				},
				{
					Event:  instance.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPRelationalAdded,
				},
				{
					Event:  instance.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPRelationalChanged,
				},
				{
					Event:  instance.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPRelationalAdded,
				},
				{
					Event:  instance.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPRelationalChanged,
				},
				{
					Event:  instance.GitLabIDPAddedEventType,
					Reduce: p.reduceGitLabIDPRelationalAdded,
				},
				{
					Event:  instance.GitLabIDPChangedEventType,
					Reduce: p.reduceGitLabIDPRelationalChanged,
				},
				{
					Event:  instance.GitLabSelfHostedIDPAddedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPRelationalAdded,
				},
				{
					Event:  instance.GitLabSelfHostedIDPChangedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPRelationalChanged,
				},
				{
					Event:  instance.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPRelationalAdded,
				},
				{
					Event:  instance.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPRelationalChanged,
				},
				{
					Event:  instance.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  instance.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  instance.AppleIDPAddedEventType,
					Reduce: p.reduceAppleIDPAdded,
				},
				{
					Event:  instance.AppleIDPChangedEventType,
					Reduce: p.reduceAppleIDPChanged,
				},
				{
					Event:  instance.SAMLIDPAddedEventType,
					Reduce: p.reduceSAMLIDPAdded,
				},
				{
					Event:  instance.SAMLIDPChangedEventType,
					Reduce: p.reduceSAMLIDPChanged,
				},
				{
					Event:  instance.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.IDPConfigAddedEventType,
					Reduce: p.reduceIDPRelationalAdded,
				},
				{
					Event:  org.IDPConfigChangedEventType,
					Reduce: p.reduceIDPRelationalChanged,
				},
				{
					Event:  org.IDPConfigDeactivatedEventType,
					Reduce: p.reduceIDRelationalPDeactivated,
				},
				{
					Event:  org.IDPConfigReactivatedEventType,
					Reduce: p.reduceIDPRelationalReactivated,
				},
				{
					Event:  org.IDPConfigRemovedEventType,
					Reduce: p.reduceIDPRelationalRemoved,
				},
				{
					Event:  org.IDPOIDCConfigAddedEventType,
					Reduce: p.reduceOIDCRelationalConfigAdded,
				},
				{
					Event:  org.IDPOIDCConfigChangedEventType,
					Reduce: p.reduceOIDCRelationalConfigChanged,
				},
				{
					Event:  org.IDPJWTConfigAddedEventType,
					Reduce: p.reduceJWTRelationalConfigAdded,
				},
				{
					Event:  org.IDPJWTConfigChangedEventType,
					Reduce: p.reduceJWTRelationalConfigChanged,
				},
				{
					Event:  org.OAuthIDPAddedEventType,
					Reduce: p.reduceOAuthIDPRelationalAdded,
				},
				{
					Event:  org.OAuthIDPChangedEventType,
					Reduce: p.reduceOAuthIDPRelationalChanged,
				},
				{
					Event:  org.OIDCIDPAddedEventType,
					Reduce: p.reduceOIDCIDPRelationalAdded,
				},
				{
					Event:  org.OIDCIDPChangedEventType,
					Reduce: p.reduceOIDCIDPRelationalChanged,
				},
				{
					Event:  org.OIDCIDPMigratedAzureADEventType,
					Reduce: p.reduceOIDCIDPRelationalMigratedAzureAD,
				},
				{
					Event:  org.OIDCIDPMigratedGoogleEventType,
					Reduce: p.reduceOIDCIDPRelationalMigratedGoogle,
				},
				{
					Event:  org.JWTIDPAddedEventType,
					Reduce: p.reduceJWTIDPRelationalAdded,
				},
				{
					Event:  org.JWTIDPChangedEventType,
					Reduce: p.reduceJWTIDPRelationalChanged,
				},
				{
					Event:  org.AzureADIDPAddedEventType,
					Reduce: p.reduceAzureADIDPRelationalAdded,
				},
				{
					Event:  org.AzureADIDPChangedEventType,
					Reduce: p.reduceAzureADIDPRelationalChanged,
				},
				{
					Event:  org.GitHubIDPAddedEventType,
					Reduce: p.reduceGitHubIDPRelationalAdded,
				},
				{
					Event:  org.GitHubIDPChangedEventType,
					Reduce: p.reduceGitHubIDPRelationalChanged,
				},
				{
					Event:  org.GitHubEnterpriseIDPAddedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPRelationalAdded,
				},
				{
					Event:  org.GitHubEnterpriseIDPChangedEventType,
					Reduce: p.reduceGitHubEnterpriseIDPRelationalChanged,
				},
				{
					Event:  org.GitLabIDPAddedEventType,
					Reduce: p.reduceGitLabIDPRelationalAdded,
				},
				{
					Event:  org.GitLabIDPChangedEventType,
					Reduce: p.reduceGitLabIDPRelationalChanged,
				},
				{
					Event:  org.GitLabSelfHostedIDPAddedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPRelationalAdded,
				},
				{
					Event:  org.GitLabSelfHostedIDPChangedEventType,
					Reduce: p.reduceGitLabSelfHostedIDPRelationalChanged,
				},
				{
					Event:  org.GoogleIDPAddedEventType,
					Reduce: p.reduceGoogleIDPRelationalAdded,
				},
				{
					Event:  org.GoogleIDPChangedEventType,
					Reduce: p.reduceGoogleIDPRelationalChanged,
				},
				{
					Event:  org.LDAPIDPAddedEventType,
					Reduce: p.reduceLDAPIDPAdded,
				},
				{
					Event:  org.LDAPIDPChangedEventType,
					Reduce: p.reduceLDAPIDPChanged,
				},
				{
					Event:  org.AppleIDPAddedEventType,
					Reduce: p.reduceAppleIDPAdded,
				},
				{
					Event:  org.AppleIDPChangedEventType,
					Reduce: p.reduceAppleIDPChanged,
				},
				{
					Event:  org.SAMLIDPAddedEventType,
					Reduce: p.reduceSAMLIDPAdded,
				},
				{
					Event:  org.SAMLIDPChangedEventType,
					Reduce: p.reduceSAMLIDPChanged,
				},
				{
					Event:  org.IDPRemovedEventType,
					Reduce: p.reduceIDPRemoved,
				},
			},
		},
	}
}

func (p *idpTemplateRelationalProjection) reduceIDPRelationalAdded(event eventstore.Event) (*handler.Statement, error) {
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
		return nil
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

func (p *idpTemplateRelationalProjection) reduceIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-9sX8h", "reduce.wrong.db.pool %T", ex)
		}

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, event.Aggregate().InstanceID, idpEvent.ConfigID, orgID), changes...)
		return err
	}), nil
}

func (p *idpTemplateRelationalProjection) reduceIDPDeactivated(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceIDPReactivated(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceIDPRemoved(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-PSj7F", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, event.Aggregate().InstanceID, idpID, orgID))
		return err
	}), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-5cvzY", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		_, err = repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.IDPConfigID, orgID),
			repo.SetPayload(string(payloadJSON)),
			database.NewChange(repo.TypeColumn(), domain.IDPTypeOIDC),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-sh6Lp", "unable to cast to tx executer")
		}

		idpRepo := repository.IDProviderRepository()
		oidc, err := idpRepo.GetOIDC(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.IDPConfigID, orgId)))
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
			database.NewChange(idpRepo.TypeColumn(), domain.IDPTypeOIDC),
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

func (p *idpTemplateRelationalProjection) reduceJWTConfigAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-tJQ8V", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		_, err = repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.IDPConfigID, orgID),
			repo.SetPayload(string(payloadJSON)),
			database.NewChange(repo.TypeColumn(), domain.IDPTypeJWT),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil
}

func (p *idpTemplateRelationalProjection) reduceJWTConfigChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInternal(nil, "HANDL-sh6Lp", "unable to cast to tx executer")
		}

		idpRepo := repository.IDProviderRepository()
		jwt, err := idpRepo.GetJWT(ctx, v3_sql.SQLTx(tx), database.WithCondition(idpScopedCondition(idpRepo, idpEvent.Agg.InstanceID, idpEvent.IDPConfigID, orgId)))
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

		payloadJSON, err := json.Marshal(jwt.JWT)
		if err != nil {
			return err
		}

		_, err = idpRepo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(idpRepo, idpEvent.Aggregate().InstanceID, idpEvent.IDPConfigID, orgId),
			idpRepo.SetPayload(string(payloadJSON)),
			database.NewChange(idpRepo.TypeColumn(), domain.IDPTypeJWT),
			idpRepo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		)
		return err
	}), nil
}

func (p *idpTemplateRelationalProjection) reduceOAuthIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceOAuthIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceOIDCIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceOIDCIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceOIDCIDPMigratedAzureAD(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-mj7LQ", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		changes := database.Changes{
			repo.SetName(idpEvent.Name),
			database.NewChange(repo.TypeColumn(), domain.IDPTypeAzure),
			repo.SetAllowCreation(idpEvent.IsCreationAllowed),
			repo.SetAllowLinking(idpEvent.IsLinkingAllowed),
			repo.SetAllowAutoCreation(idpEvent.IsAutoCreation),
			repo.SetAllowAutoUpdate(idpEvent.IsAutoUpdate),
			repo.SetLinkingField(mapAutoLinkingField(idpEvent.AutoLinkingOption)),
			repo.SetPayload(string(payloadJSON)),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		}

		_, err = repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.ID, orgID), changes...)
		return err
	}), nil
}

func (p *idpTemplateRelationalProjection) reduceOIDCIDPMigratedGoogle(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-HDqk9", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.IDProviderRepository()
		changes := database.Changes{
			repo.SetName(idpEvent.Name),
			database.NewChange(repo.TypeColumn(), domain.IDPTypeGoogle),
			repo.SetAllowCreation(idpEvent.IsCreationAllowed),
			repo.SetAllowLinking(idpEvent.IsLinkingAllowed),
			repo.SetAllowAutoCreation(idpEvent.IsAutoCreation),
			repo.SetAllowAutoUpdate(idpEvent.IsAutoUpdate),
			repo.SetLinkingField(mapAutoLinkingField(idpEvent.AutoLinkingOption)),
			repo.SetPayload(string(payloadJSON)),
			repo.SetUpdatedAt(gu.Ptr(event.CreatedAt())),
		}

		_, err = repo.Update(ctx, v3_sql.SQLTx(tx), idpScopedCondition(repo, idpEvent.Aggregate().InstanceID, idpEvent.ID, orgID), changes...)
		return err
	}), nil
}

func (p *idpTemplateRelationalProjection) reduceJWTIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceJWTIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceAzureADIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceAzureADIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGitHubIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGitHubIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGitHubEnterpriseIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGitHubEnterpriseIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGitLabIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGitLabIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGitLabSelfHostedIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGitLabSelfHostedIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGoogleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceGoogleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceLDAPIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceLDAPIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceAppleIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceAppleIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceSAMLIDPAdded(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceSAMLIDPChanged(event eventstore.Event) (*handler.Statement, error) {
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

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, projectionName string) error {
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

func (p *idpTemplateRelationalProjection) reduceIDPChangedTemplateColumns(repo domain.IDProviderRepository, name *string, optionChanges idp.OptionChanges) database.Changes {
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
		changes = append(changes, repo.SetLinkingField(mapAutoLinkingField(*optionChanges.AutoLinkingOption)))
	}

	return changes
}

func (p *idpTemplateRelationalProjection) reduceOAuthIDPChangedColumns(payload *domain.OAuth, idpEvent *idp.OAuthIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceOIDCIDPChangedColumns(payload *domain.OIDC, idpEvent *idp.OIDCIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceJWTIDPChangedColumns(payload *domain.JWT, idpEvent *idp.JWTIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceAzureADIDPChangedColumns(payload *domain.Azure, idpEvent *idp.AzureADIDPChangedEvent) (bool, error) {
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

func (p *idpTemplateRelationalProjection) reduceGitHubIDPChangedColumns(payload *domain.Github, idpEvent *idp.GitHubIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceGitHubEnterpriseIDPChangedColumns(payload *domain.GithubEnterprise, idpEvent *idp.GitHubEnterpriseIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceGitLabIDPChangedColumns(payload *domain.Gitlab, idpEvent *idp.GitLabIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceGitLabSelfHostedIDPChangedColumns(payload *domain.GitlabSelfHosted, idpEvent *idp.GitLabSelfHostedIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceGoogleIDPChangedColumns(payload *domain.Google, idpEvent *idp.GoogleIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceLDAPIDPChangedColumns(payload *domain.LDAP, idpEvent *idp.LDAPIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceAppleIDPChangedColumns(payload *domain.Apple, idpEvent *idp.AppleIDPChangedEvent) bool {
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

func (p *idpTemplateRelationalProjection) reduceSAMLIDPChangedColumns(payload *domain.SAML, idpEvent *idp.SAMLIDPChangedEvent) bool {
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
