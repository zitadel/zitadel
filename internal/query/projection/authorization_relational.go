package projection

import (
	"context"
	"database/sql"

	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	AuthorizationRelationalProjectionTable = "zitadel.authorizations"
)

type authorizationRelationalProjection struct{}

func (*authorizationRelationalProjection) Name() string {
	return AuthorizationRelationalProjectionTable
}

func newAuthorizationRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(authorizationRelationalProjection))
}

func (a *authorizationRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: usergrant.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  usergrant.UserGrantAddedType,
					Reduce: a.reduceUserGrantAdded,
				},
				{
					Event:  usergrant.UserGrantChangedType,
					Reduce: a.reduceUserGrantChanged,
				},
				{
					Event:  usergrant.UserGrantCascadeChangedType,
					Reduce: a.reduceUserGrantCascadeChanged,
				},
				{
					Event:  usergrant.UserGrantCascadeRemovedType,
					Reduce: a.reduceUserGrantCascadeRemoved,
				},
				{
					Event:  usergrant.UserGrantRemovedType,
					Reduce: a.reduceUserGrantRemoved,
				},
				{
					Event:  usergrant.UserGrantDeactivatedType,
					Reduce: a.reduceUserGrantDeactivated,
				},
				{
					Event:  usergrant.UserGrantReactivatedType,
					Reduce: a.reduceUserGrantReactivated,
				},
			},
		},
	}
}

func (a *authorizationRelationalProjection) reduceUserGrantAdded(event eventstore.Event) (*handler.Statement, error) {
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
		return repo.Create(ctx, v3_sql.SQLTx(tx), &repoDomain.Authorization{
			InstanceID: e.Aggregate().InstanceID,
			ID:         e.Aggregate().ID,
			UserID:     e.UserID,
			ProjectID:  e.ProjectID,
			GrantID:    e.ProjectGrantID,
			Roles:      e.RoleKeys,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
			State:      repoDomain.AuthorizationStateActive,
		})
	}), nil
}

func (a *authorizationRelationalProjection) reduceUserGrantChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-Kdr7KP", "reduce.wrong.event.type %s", usergrant.UserGrantChangedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-sBwGEW", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.ID),
			repo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (a *authorizationRelationalProjection) reduceUserGrantCascadeChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantCascadeChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-dFtRl7", "reduce.wrong.event.type %s", usergrant.UserGrantCascadeChangedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55M", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.ID),
			repo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (a *authorizationRelationalProjection) reduceUserGrantCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantCascadeRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl8", "reduce.wrong.event.type %s", usergrant.UserGrantCascadeRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55N", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.ID),
		)
		return err
	}), nil
}

func (a *authorizationRelationalProjection) reduceUserGrantRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl9", "reduce.wrong.event.type %s", usergrant.UserGrantRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55O", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.ID),
		)
		return err
	}), nil
}

func (a *authorizationRelationalProjection) reduceUserGrantDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl0", "reduce.wrong.event.type %s", usergrant.UserGrantDeactivatedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55P", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.ID),
			repo.SetUpdatedAt(e.CreationDate()),
			repo.SetState(repoDomain.AuthorizationStateInactive),
		)
		return err
	}), nil
}

func (a *authorizationRelationalProjection) reduceUserGrantReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl1", "reduce.wrong.event.type %s", usergrant.UserGrantReactivatedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55Q", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.ID),
			repo.SetUpdatedAt(e.CreationDate()),
			repo.SetState(repoDomain.AuthorizationStateActive),
		)
		return err
	}), nil
}
