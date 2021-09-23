package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/project"
)

type ProjectRoleProjection struct {
	crdb.StatementHandler
}

const ProjectRoleProjectionTable = "zitadel.projections.project_roles"

func NewProjectRoleProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectRoleProjection {
	p := &ProjectRoleProjection{}
	config.ProjectionName = ProjectRoleProjectionTable
	config.Reducers = p.reducers()
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
			},
		},
	}
}

const (
	ProjectRoleProjectIDCol     = "project_id"
	ProjectRoleKeyCol           = "role_key"
	ProjectRoleCreationDateCol  = "creation_date"
	ProjectRoleChangeDateCol    = "change_date"
	ProjectRoleResourceOwnerCol = "resource_owner"
	ProjectRoleSequenceCol      = "sequence"
	ProjectRoleDisplayNameCol   = "display_name"
	ProjectRoleGroupNameCol     = "group_name"
	ProjectRoleCreatorCol       = "creator_id"
)

func (p *ProjectRoleProjection) reduceProjectRoleAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.RoleAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Fmre5", "seq", event.Sequence(), "expectedType", project.RoleAddedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-g92Fg", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectRoleKeyCol, e.Key),
			handler.NewCol(ProjectRoleProjectIDCol, e.Aggregate().ID),
			handler.NewCol(ProjectRoleCreationDateCol, e.CreationDate()),
			handler.NewCol(ProjectRoleChangeDateCol, e.CreationDate()),
			handler.NewCol(ProjectRoleResourceOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectRoleSequenceCol, e.Sequence()),
			handler.NewCol(ProjectRoleDisplayNameCol, e.DisplayName),
			handler.NewCol(ProjectRoleGroupNameCol, e.Group),
			handler.NewCol(ProjectRoleCreatorCol, e.EditorUser()),
		},
	), nil
}

func (p *ProjectRoleProjection) reduceProjectRoleChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.RoleChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-M0fwg", "seq", event.Sequence(), "expectedType", project.GrantChangedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-sM0f", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectChangeDateCol, e.CreationDate()),
			handler.NewCol(ProjectRoleSequenceCol, e.Sequence()),
			handler.NewCol(ProjectRoleDisplayNameCol, e.DisplayName),
			handler.NewCol(ProjectRoleGroupNameCol, e.Group),
		},
		[]handler.Condition{
			handler.NewCond(ProjectRoleKeyCol, e.Key),
			handler.NewCond(ProjectRoleProjectIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectRoleProjection) reduceProjectRoleRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.RoleRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-MlokF", "seq", event.Sequence(), "expectedType", project.GrantRemovedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-L0fJf", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectRoleKeyCol, e.Key),
			handler.NewCond(ProjectRoleProjectIDCol, e.Aggregate().ID),
		},
	), nil
}
