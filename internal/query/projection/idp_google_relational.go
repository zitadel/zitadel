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
