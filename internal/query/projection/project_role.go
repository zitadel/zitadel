package projection

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/project"
)

const (
	ProjectRoleProjectionTable = "projections.project_roles"

	ProjectRoleColumnProjectID     = "project_id"
	ProjectRoleColumnKey           = "role_key"
	ProjectRoleColumnCreationDate  = "creation_date"
	ProjectRoleColumnChangeDate    = "change_date"
	ProjectRoleColumnSequence      = "sequence"
	ProjectRoleColumnResourceOwner = "resource_owner"
	ProjectRoleColumnDisplayName   = "display_name"
	ProjectRoleColumnGroupName     = "group_name"
	ProjectRoleColumnCreator       = "creator_id" //TODO: necessary?
)

type ProjectRoleProjection struct {
	crdb.StatementHandler
}

func NewProjectRoleProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectRoleProjection {
	p := new(ProjectRoleProjection)
	config.ProjectionName = ProjectRoleProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(ProjectRoleColumnProjectID, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectRoleColumnKey, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectRoleColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(ProjectRoleColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(ProjectRoleColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(ProjectRoleColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectRoleColumnDisplayName, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectRoleColumnGroupName, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectRoleColumnCreator, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(ProjectRoleColumnProjectID, ProjectRoleColumnKey),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *ProjectRoleProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
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
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
	}
}

func (p *ProjectRoleProjection) reduceProjectRoleAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-g92Fg", "reduce.wrong.event.type %s", project.RoleAddedType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectRoleColumnKey, e.Key),
			handler.NewCol(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCol(ProjectRoleColumnCreationDate, e.CreationDate()),
			handler.NewCol(ProjectRoleColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectRoleColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectRoleColumnSequence, e.Sequence()),
			handler.NewCol(ProjectRoleColumnDisplayName, e.DisplayName),
			handler.NewCol(ProjectRoleColumnGroupName, e.Group),
			handler.NewCol(ProjectRoleColumnCreator, e.EditorUser()),
		},
	), nil
}

func (p *ProjectRoleProjection) reduceProjectRoleChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-sM0f", "reduce.wrong.event.type %s", project.GrantChangedType)
	}
	if e.DisplayName == nil && e.Group == nil {
		return crdb.NewNoOpStatement(e), nil
	}
	columns := make([]handler.Column, 0, 7)
	columns = append(columns, handler.NewCol(ProjectRoleColumnChangeDate, e.CreationDate()),
		handler.NewCol(ProjectRoleColumnSequence, e.Sequence()))
	if e.DisplayName != nil {
		columns = append(columns, handler.NewCol(ProjectRoleColumnDisplayName, *e.DisplayName))
	}
	if e.Group != nil {
		columns = append(columns, handler.NewCol(ProjectRoleColumnGroupName, *e.Group))
	}
	return crdb.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnKey, e.Key),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectRoleProjection) reduceProjectRoleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-L0fJf", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnKey, e.Key),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectRoleProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-l0geG", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
		},
	), nil
}
