package projection

import (
	"context"
	"database/sql"

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
