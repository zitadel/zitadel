package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

const (
	ProjectRoleProjectionTable = "projections.project_roles3"

	ProjectRoleColumnProjectID     = "project_id"
	ProjectRoleColumnKey           = "role_key"
	ProjectRoleColumnCreationDate  = "creation_date"
	ProjectRoleColumnChangeDate    = "change_date"
	ProjectRoleColumnSequence      = "sequence"
	ProjectRoleColumnResourceOwner = "resource_owner"
	ProjectRoleColumnInstanceID    = "instance_id"
	ProjectRoleColumnDisplayName   = "display_name"
	ProjectRoleColumnGroupName     = "group_name"
	ProjectRoleColumnOwnerRemoved  = "owner_removed"
)

type projectRoleProjection struct{}

func newProjectRoleProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(projectRoleProjection))
}

func (*projectRoleProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(ProjectRoleColumnProjectID, handler.ColumnTypeText),
			handler.NewColumn(ProjectRoleColumnKey, handler.ColumnTypeText),
			handler.NewColumn(ProjectRoleColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(ProjectRoleColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(ProjectRoleColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(ProjectRoleColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(ProjectRoleColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(ProjectRoleColumnDisplayName, handler.ColumnTypeText),
			handler.NewColumn(ProjectRoleColumnGroupName, handler.ColumnTypeText),
			handler.NewColumn(ProjectRoleColumnOwnerRemoved, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(ProjectRoleColumnInstanceID, ProjectRoleColumnProjectID, ProjectRoleColumnKey),
			handler.WithIndex(handler.NewIndex("owner_removed", []string{ProjectRoleColumnOwnerRemoved})),
		),
	)
}

func (*projectRoleProjection) Name() string {
	return ProjectRoleProjectionTable
}

func (p *projectRoleProjection) Reducers() []handler.AggregateReducer {
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
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(ProjectRoleColumnInstanceID),
				},
			},
		},
	}
}

func (p *projectRoleProjection) reduceProjectRoleAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-g92Fg", "reduce.wrong.event.type %s", project.RoleAddedType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectRoleColumnKey, e.Key),
			handler.NewCol(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCol(ProjectRoleColumnCreationDate, e.CreationDate()),
			handler.NewCol(ProjectRoleColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectRoleColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectRoleColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(ProjectRoleColumnSequence, e.Sequence()),
			handler.NewCol(ProjectRoleColumnDisplayName, e.DisplayName),
			handler.NewCol(ProjectRoleColumnGroupName, e.Group),
		},
	), nil
}

func (p *projectRoleProjection) reduceProjectRoleChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-sM0f", "reduce.wrong.event.type %s", project.GrantChangedType)
	}
	if e.DisplayName == nil && e.Group == nil {
		return handler.NewNoOpStatement(e), nil
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
	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnKey, e.Key),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectRoleColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectRoleProjection) reduceProjectRoleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-L0fJf", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnKey, e.Key),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectRoleColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectRoleProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-l0geG", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectRoleColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectRoleProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-3XrHY", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ProjectRoleColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
