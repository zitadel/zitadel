package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type projectRoleProjection struct {
	crdb.StatementHandler
}

const ProjectRoleProjectionTable = "zitadel.projections.project_roles"

func newProjectRoleProjection(ctx context.Context, config crdb.StatementHandlerConfig) *projectRoleProjection {
	p := &projectRoleProjection{}
	config.ProjectionName = ProjectRoleProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *projectRoleProjection) reducers() []handler.AggregateReducer {
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

const (
	ProjectRoleColumnProjectID     = "project_id"
	ProjectRoleColumnKey           = "role_key"
	ProjectRoleColumnCreationDate  = "creation_date"
	ProjectRoleColumnChangeDate    = "change_date"
	ProjectRoleColumnResourceOwner = "resource_owner"
	ProjectRoleColumnSequence      = "sequence"
	ProjectRoleColumnDisplayName   = "display_name"
	ProjectRoleColumnGroupName     = "group_name"
	ProjectRoleColumnCreator       = "creator_id"
)

func (p *projectRoleProjection) reduceProjectRoleAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Fmre5", "seq", event.Sequence(), "expectedType", project.RoleAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-g92Fg", "reduce.wrong.event.type")
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

func (p *projectRoleProjection) reduceProjectRoleChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-M0fwg", "seq", event.Sequence(), "expectedType", project.GrantChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-sM0f", "reduce.wrong.event.type")
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

func (p *projectRoleProjection) reduceProjectRoleRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-MlokF", "seq", event.Sequence(), "expectedType", project.GrantRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-L0fJf", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnKey, e.Key),
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *projectRoleProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-hm90R", "seq", event.Sequence(), "expectedType", project.ProjectRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-l0geG", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleColumnProjectID, e.Aggregate().ID),
		},
	), nil
}
