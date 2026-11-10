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
