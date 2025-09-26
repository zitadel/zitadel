package projection

import (
	"context"

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
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(InstanceID, e.Aggregate().InstanceID),
			handler.NewCol(OrganizationID, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCol(CreatedAt, e.CreationDate()),
			handler.NewCol(UpdatedAt, e.CreationDate()),
			handler.NewCol(ProjectRoleRelationalKeyCol, e.Key),
			handler.NewCol(ProjectRoleRelationalDisplayNameCol, e.DisplayName),
			handler.NewCol(ProjectRoleRelationalRoleGroupCol, e.Group),
		},
	), nil
}

func (p *projectRoleRelationalProjection) reduceProjectRoleChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-jie5J", "reduce.wrong.event.type %s", project.GrantChangedType)
	}
	if e.DisplayName == nil && e.Group == nil {
		return handler.NewNoOpStatement(e), nil
	}
	columns := make([]handler.Column, 0, 6)
	columns = append(columns,
		handler.NewCol(UpdatedAt, e.CreationDate()),
	)
	if e.DisplayName != nil {
		columns = append(columns, handler.NewCol(ProjectRoleRelationalDisplayNameCol, *e.DisplayName))
	}
	if e.Group != nil {
		columns = append(columns, handler.NewCol(ProjectRoleRelationalRoleGroupCol, *e.Group))
	}
	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(InstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectRoleRelationalKeyCol, e.Key),
		},
	), nil
}

func (p *projectRoleRelationalProjection) reduceProjectRoleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-euf0U", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(InstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectRoleRelationalKeyCol, e.Key),
		},
	), nil
}
