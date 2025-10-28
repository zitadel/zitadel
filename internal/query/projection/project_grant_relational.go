package projection

import (
	"context"
	"database/sql"

	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type projectGrantRelationalProjection struct{}

func (*projectGrantRelationalProjection) Name() string {
	return "zitadel.project_grants"
}

func newProjectGrantRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(projectGrantRelationalProjection))
}

func (p *projectGrantRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.GrantAddedType,
					Reduce: p.reduceProjectGrantAdded,
				},
				{
					Event:  project.GrantChangedType,
					Reduce: p.reduceProjectGrantChanged,
				},
				{
					Event:  project.GrantCascadeChangedType,
					Reduce: p.reduceProjectGrantCascadeChanged,
				},
				{
					Event:  project.GrantDeactivatedType,
					Reduce: p.reduceProjectGrantDeactivated,
				},
				{
					Event:  project.GrantReactivatedType,
					Reduce: p.reduceProjectGrantReactivated,
				},
				{
					Event:  project.GrantRemovedType,
					Reduce: p.reduceProjectGrantRemoved,
				},
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
	}
}

func (p *projectGrantRelationalProjection) reduceProjectGrantAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.GrantAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.ProjectGrantRepository()
		if err := repo.Create(ctx, v3_sql.SQLTx(tx), &repoDomain.ProjectGrant{
			InstanceID:             e.Aggregate().InstanceID,
			ID:                     e.GrantID,
			ProjectID:              e.Aggregate().ID,
			GrantedOrganizationID:  e.GrantedOrgID,
			GrantingOrganizationID: e.Aggregate().ResourceOwner,
			CreatedAt:              e.CreationDate(),
			UpdatedAt:              e.CreationDate(),
			State:                  repoDomain.ProjectGrantStateActive,
		}); err != nil {
			return err
		}
		if len(e.RoleKeys) > 0 {
			roleRepo := repo.Role()
			for _, key := range e.RoleKeys {
				if err := roleRepo.Add(ctx, v3_sql.SQLTx(tx), &repoDomain.ProjectGrantRole{
					InstanceID:   e.Aggregate().InstanceID,
					GrantID:      e.GrantID,
					ProjectID:    e.Aggregate().ID,
					ProjectOrgID: e.Aggregate().ResourceOwner,
					CreatedAt:    e.CreationDate(),
					Key:          key,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	}), nil
}

func (p *projectGrantRelationalProjection) reduceProjectGrantChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.GrantChangedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.ProjectGrantRepository().Role()
		if _, err := repo.Remove(ctx, v3_sql.SQLTx(tx),
			database.And(
				repo.InstanceIDCondition(e.Aggregate().InstanceID),
				repo.GrantIDCondition(e.GrantID),
			),
		); err != nil {
			return err
		}
		if len(e.RoleKeys) > 0 {
			for _, key := range e.RoleKeys {
				if err := repo.Add(ctx, v3_sql.SQLTx(tx), &repoDomain.ProjectGrantRole{
					InstanceID:   e.Aggregate().InstanceID,
					GrantID:      e.GrantID,
					ProjectID:    e.Aggregate().ID,
					ProjectOrgID: e.Aggregate().ResourceOwner,
					CreatedAt:    e.CreationDate(),
					Key:          key,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	}), nil
}

func (p *projectGrantRelationalProjection) reduceProjectGrantCascadeChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantCascadeChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.GrantCascadeChangedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.ProjectGrantRepository().Role()
		if _, err := repo.Remove(ctx, v3_sql.SQLTx(tx),
			database.And(
				repo.InstanceIDCondition(e.Aggregate().InstanceID),
				repo.GrantIDCondition(e.GrantID),
			),
		); err != nil {
			return err
		}
		if len(e.RoleKeys) > 0 {
			for _, key := range e.RoleKeys {
				if err := repo.Add(ctx, v3_sql.SQLTx(tx), &repoDomain.ProjectGrantRole{
					InstanceID:   e.Aggregate().InstanceID,
					GrantID:      e.GrantID,
					ProjectID:    e.Aggregate().ID,
					ProjectOrgID: e.Aggregate().ResourceOwner,
					CreatedAt:    e.CreationDate(),
					Key:          key,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	}), nil
}

func (p *projectGrantRelationalProjection) reduceProjectGrantDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantDeactivateEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.GrantDeactivatedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		repo := repository.ProjectGrantRepository()
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.GrantID),
			repo.SetUpdatedAt(e.CreationDate()),
			repo.SetState(repoDomain.ProjectGrantStateInactive),
		)
		return err
	}), nil
}

func (p *projectGrantRelationalProjection) reduceProjectGrantReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-oof4U", "reduce.wrong.event.type %s", project.GrantReactivatedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		repo := repository.ProjectGrantRepository()
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetUpdatedAt(e.CreationDate()),
			repo.SetState(repoDomain.ProjectGrantStateActive),
		)
		return err
	}), nil
}

func (p *projectGrantRelationalProjection) reduceProjectGrantRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Xae7w", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		repo := repository.ProjectGrantRepository()
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.GrantID),
		)
		return err
	}), nil
}

func (p *projectGrantRelationalProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Xae7w", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		repo := repository.ProjectGrantRepository()
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			database.And(
				repo.InstanceIDCondition(e.Aggregate().InstanceID),
				repo.ProjectIDCondition(e.Aggregate().ID),
			),
		)
		return err
	}), nil
}
