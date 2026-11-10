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
