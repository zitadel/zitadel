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
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/idp"
	"github.com/zitadel/zitadel/internal/repository/idpconfig"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

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
