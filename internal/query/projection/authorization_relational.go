package projection

import (
	"context"
	"database/sql"

	"github.com/muhlemmer/gu"

	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceAuthorizationAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-an0k9a", "reduce.wrong.event.type %s", usergrant.UserGrantAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FrDVPk", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.AuthorizationRepository()
		var grantID *string
		if e.ProjectGrantID != "" {
			grantID = gu.Ptr(e.ProjectGrantID)
		}
		return repo.Create(ctx, v3_sql.SQLTx(tx), &repoDomain.Authorization{
			InstanceID: e.Aggregate().InstanceID,
			ID:         e.Aggregate().ID,
			UserID:     e.UserID,
			ProjectID:  e.ProjectID,
			GrantID:    grantID,
			Roles:      e.RoleKeys,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
			State:      repoDomain.AuthorizationStateActive,
		})
	}), nil
}

func (p *relationalTablesProjection) reduceAuthorizationChanged(event eventstore.Event) (*handler.Statement, error) {
	var roles []string
	switch e := event.(type) {
	case *usergrant.UserGrantChangedEvent:
		roles = e.RoleKeys
	case *usergrant.UserGrantCascadeChangedEvent:
		roles = e.RoleKeys
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-kg1ryL", "reduce.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantChangedType, usergrant.UserGrantCascadeChangedType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-sBwGEW", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(event.Aggregate().InstanceID, event.Aggregate().ID),
			roles,
			repo.SetUpdatedAt(event.CreatedAt()),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceAuthorizationRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *usergrant.UserGrantRemovedEvent, *usergrant.UserGrantCascadeRemovedEvent:
		// ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl9", "reduce.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantRemovedType, usergrant.UserGrantCascadeRemovedType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55O", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(event.Aggregate().InstanceID, event.Aggregate().ID),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceAuthorizationDeactivated(event eventstore.Event) (*handler.Statement, error) {
	_, ok := event.(*usergrant.UserGrantDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl0", "reduce.wrong.event.type %s", usergrant.UserGrantDeactivatedType)
	}
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55P", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(event.Aggregate().InstanceID, event.Aggregate().ID),
			nil,
			repo.SetUpdatedAt(event.CreatedAt()),
			repo.SetState(repoDomain.AuthorizationStateInactive),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceAuthorizationReactivated(event eventstore.Event) (*handler.Statement, error) {
	_, ok := event.(*usergrant.UserGrantReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl1", "reduce.wrong.event.type %s", usergrant.UserGrantReactivatedType)
	}
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55Q", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(event.Aggregate().InstanceID, event.Aggregate().ID),
			nil,
			repo.SetUpdatedAt(event.CreatedAt()),
			repo.SetState(repoDomain.AuthorizationStateActive),
		)
		return err
	}), nil
}
