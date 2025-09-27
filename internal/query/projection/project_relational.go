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

type projectRelationalProjection struct{}

func (*projectRelationalProjection) Name() string {
	return "zitadel.projects"
}

func newProjectRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(projectRelationalProjection))
}

func (p *projectRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.ProjectAddedType,
					Reduce: p.reduceProjectAdded,
				},
				{
					Event:  project.ProjectChangedType,
					Reduce: p.reduceProjectChanged,
				},
				{
					Event:  project.ProjectDeactivatedType,
					Reduce: p.reduceProjectDeactivated,
				},
				{
					Event:  project.ProjectReactivatedType,
					Reduce: p.reduceProjectReactivated,
				},
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
	}
}

func (p *projectRelationalProjection) reduceProjectAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.ProjectAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.ProjectRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &repoDomain.Project{
			InstanceID:               e.Aggregate().InstanceID,
			OrganizationID:           e.Aggregate().ResourceOwner,
			ID:                       e.Aggregate().ID,
			CreatedAt:                e.CreationDate(),
			UpdatedAt:                e.CreationDate(),
			Name:                     e.Name,
			State:                    repoDomain.ProjectStateActive,
			ShouldAssertRole:         e.ProjectRoleAssertion,
			IsAuthorizationRequired:  e.ProjectRoleCheck,
			IsProjectAccessRequired:  e.HasProjectCheck,
			UsedLabelingSettingOwner: int16(e.PrivateLabelingSetting),
		})
	}), nil
}

func (p *projectRelationalProjection) reduceProjectChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectChangeEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.ProjectChangedType)
	}
	if e.Name == nil && e.HasProjectCheck == nil && e.ProjectRoleAssertion == nil && e.ProjectRoleCheck == nil && e.PrivateLabelingSetting == nil {
		return handler.NewNoOpStatement(e), nil
	}

	repo := repository.ProjectRepository()
	changes := make([]database.Change, 0, 6)
	changes = append(changes, repo.SetUpdatedAt(e.CreationDate()))
	if e.Name != nil {
		changes = append(changes, repo.SetName(*e.Name))
	}
	if e.ProjectRoleAssertion != nil {
		changes = append(changes, repo.SetShouldAssertRole(*e.ProjectRoleAssertion))
	}
	if e.ProjectRoleCheck != nil {
		changes = append(changes, repo.SetIsAuthorizationRequired(*e.ProjectRoleCheck))
	}
	if e.HasProjectCheck != nil {
		changes = append(changes, repo.SetIsProjectAccessRequired(*e.HasProjectCheck))
	}
	if e.PrivateLabelingSetting != nil {
		changes = append(changes, repo.SetUsedLabelingSettingOwner(int16(*e.PrivateLabelingSetting)))
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			changes...,
		)
		return err
	}), nil
}

func (p *projectRelationalProjection) reduceProjectDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Oox5e", "reduce.wrong.event.type %s", project.ProjectDeactivatedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		repo := repository.ProjectRepository()
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetUpdatedAt(e.CreationDate()),
			repo.SetState(repoDomain.ProjectStateInactive),
		)
		return err
	}), nil
}

func (p *projectRelationalProjection) reduceProjectReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-oof4U", "reduce.wrong.event.type %s", project.ProjectDeactivatedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		repo := repository.ProjectRepository()
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetUpdatedAt(e.CreationDate()),
			repo.SetState(repoDomain.ProjectStateActive),
		)
		return err
	}), nil
}

func (p *projectRelationalProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Xae7w", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		repo := repository.ProjectRepository()
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
		)
		return err
	}), nil
}
