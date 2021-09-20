package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/project"
)

type ProjectProjection struct {
	crdb.StatementHandler
}

const (
	projectProjection = "zitadel.projections.projects"
)

func NewProjectProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectProjection {
	p := &ProjectProjection{}
	config.ProjectionName = projectProjection
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *ProjectProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
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

const (
	projectIDCol                   = "id"
	projectNameCol                 = "name"
	projectProjectRoleAssertionCol = "project_role_assertion"
	projectProjectRoleCheckCol     = "project_role_check"
	projectHasProjectCheckCol      = "has_project_check"
	projectPrivateLabelingCol      = "private_labeling_setting"
	projectCreationDateCol         = "creation_date"
	projectChangeDateCol           = "change_date"
	projectOwnerCol                = "resource_owner"
	projectCreatorCol              = "creator_id"
	projectStateCol                = "state"
	projectSequenceCol             = "sequence"
)

func (p *ProjectProjection) reduceProjectAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-MFOsd", "seq", event.Sequence(), "expectedType", project.ProjectAddedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-l000S", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectIDCol, e.Aggregate().ID),
			handler.NewCol(projectCreationDateCol, e.CreationDate()),
			handler.NewCol(projectChangeDateCol, e.CreationDate()),
			handler.NewCol(projectOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(projectSequenceCol, e.Sequence()),
			handler.NewCol(projectNameCol, e.Name),
			handler.NewCol(projectProjectRoleAssertionCol, e.ProjectRoleAssertion),
			handler.NewCol(projectProjectRoleCheckCol, e.ProjectRoleCheck),
			handler.NewCol(projectHasProjectCheckCol, e.HasProjectCheck),
			handler.NewCol(projectPrivateLabelingCol, e.PrivateLabelingSetting),
			handler.NewCol(projectStateCol, domain.ProjectStateActive),
			handler.NewCol(projectCreatorCol, e.EditorUser()),
		},
	), nil
}

func (p *ProjectProjection) reduceProjectChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectChangeEvent)
	if !ok {
		logging.LogWithFields("HANDL-dk2iF", "seq", event.Sequence(), "expected", project.ProjectChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-s00Fs", "reduce.wrong.event.type")
	}
	if e.Name == nil {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectChangeDateCol, e.CreationDate()),
			handler.NewCol(projectSequenceCol, e.Sequence()),
			handler.NewCol(projectNameCol, *e.Name),
			handler.NewCol(projectProjectRoleAssertionCol, *e.ProjectRoleAssertion),
			handler.NewCol(projectProjectRoleCheckCol, *e.ProjectRoleCheck),
			handler.NewCol(projectHasProjectCheckCol, *e.HasProjectCheck),
			handler.NewCol(projectPrivateLabelingCol, *e.PrivateLabelingSetting),
		},
		[]handler.Condition{
			handler.NewCond(projectIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectProjection) reduceProjectDeactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-8Nf2s", "seq", event.Sequence(), "expectedType", project.ProjectDeactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-LLp0f", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectChangeDateCol, e.CreationDate()),
			handler.NewCol(projectSequenceCol, e.Sequence()),
			handler.NewCol(projectStateCol, domain.ProjectStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(projectIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectProjection) reduceProjectReactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-sm99f", "seq", event.Sequence(), "expectedType", project.ProjectReactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-9J98f", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectChangeDateCol, e.CreationDate()),
			handler.NewCol(projectSequenceCol, e.Sequence()),
			handler.NewCol(projectStateCol, domain.ProjectStateActive),
		},
		[]handler.Condition{
			handler.NewCond(projectIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectProjection) reduceProjectRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-mL0sf", "seq", event.Sequence(), "expectedType", project.ProjectRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-5N9fs", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(projectIDCol, e.Aggregate().ID),
		},
	), nil
}
