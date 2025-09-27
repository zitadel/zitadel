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

const (
	ProjectRoleRelationalTable          = "zitadel.project_roles"
	ProjectRoleRelationalKeyCol         = "key"
	ProjectRoleRelationalDisplayNameCol = "display_name"
	ProjectRoleRelationalRoleGroupCol   = "role_group"
)

type projectRoleRelationalProjection struct{}

func (*projectRoleRelationalProjection) Name() string {
	return ProjectRoleRelationalTable
}

func newProjectRoleRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(projectRoleRelationalProjection))
}

func (p *projectRoleRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.RoleAddedType,
					Reduce: p.reduceProjectRoleAdded,
				},
				{
					Event:  project.RoleChangedType,
					Reduce: p.reduceProjectRoleChanged,
				},
				{
					Event:  project.RoleRemovedType,
					Reduce: p.reduceProjectRoleRemoved,
				},
			},
		},
	}
}

func (p *projectRoleRelationalProjection) reduceProjectRoleAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-uw7Oo", "reduce.wrong.event.type %s", project.RoleAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}

		// Group is optional and nullable but not a pointer in the event
		// so we need to convert the empty string to a nil pointer
		var group *string
		if e.Group != "" {
			group = &e.Group
		}

		repo := repository.ProjectRepository().Role()
		return repo.Create(ctx, v3_sql.SQLTx(tx), &repoDomain.ProjectRole{
			InstanceID:     e.Aggregate().InstanceID,
			OrganizationID: e.Aggregate().ResourceOwner,
			ProjectID:      e.Aggregate().ID,
			CreatedAt:      e.CreationDate(),
			UpdatedAt:      e.CreationDate(),
			Key:            e.Key,
			DisplayName:    e.DisplayName,
			RoleGroup:      group,
		})
	}), nil
}

func (p *projectRoleRelationalProjection) reduceProjectRoleChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-jie5J", "reduce.wrong.event.type %s", project.GrantChangedType)
	}
	if e.DisplayName == nil && e.Group == nil {
		return handler.NewNoOpStatement(e), nil
	}

	repo := repository.ProjectRepository().Role()
	changes := make([]database.Change, 0, 3)
	changes = append(changes,
		repo.SetUpdatedAt(e.CreationDate()),
	)
	if e.DisplayName != nil {
		changes = append(changes, repo.SetDisplayName(*e.DisplayName))
	}
	if e.Group != nil {
		changes = append(changes, repo.SetRoleGroup(*e.Group))
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID, e.Key),
			changes...,
		)
		return err
	}), nil
}

func (p *projectRoleRelationalProjection) reduceProjectRoleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-euf0U", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		repo := repository.ProjectRepository().Role()
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-kGokE", "reduce.wrong.db.pool %T", ex)
		}
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID, e.Key),
		)
		return err
	}), nil
}
