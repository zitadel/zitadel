package projection

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/project"
)

type projectProjection struct {
	crdb.StatementHandler
}

const (
	ProjectProjectionTable = "zitadel.projections.projects"
)

func newProjectProjection(ctx context.Context, config crdb.StatementHandlerConfig) *projectProjection {
	p := &projectProjection{}
	config.ProjectionName = ProjectProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *projectProjection) reducers() []handler.AggregateReducer {
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
	ProjectColumnID                     = "id"
	ProjectColumnName                   = "name"
	ProjectColumnProjectRoleAssertion   = "project_role_assertion"
	ProjectColumnProjectRoleCheck       = "project_role_check"
	ProjectColumnHasProjectCheck        = "has_project_check"
	ProjectColumnPrivateLabelingSetting = "private_labeling_setting"
	ProjectColumnCreationDate           = "creation_date"
	ProjectColumnChangeDate             = "change_date"
	ProjectColumnResourceOwner          = "resource_owner"
	ProjectColumnCreator                = "creator_id"
	ProjectColumnState                  = "state"
	ProjectColumnSequence               = "sequence"
)

func (p *projectProjection) reduceProjectAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-MFOsd", "seq", event.Sequence(), "expectedType", project.ProjectAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-l000S", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnID, e.Aggregate().ID),
			handler.NewCol(ProjectColumnCreationDate, e.CreationDate()),
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectColumnSequence, e.Sequence()),
			handler.NewCol(ProjectColumnName, e.Name),
			handler.NewCol(ProjectColumnProjectRoleAssertion, e.ProjectRoleAssertion),
			handler.NewCol(ProjectColumnProjectRoleCheck, e.ProjectRoleCheck),
			handler.NewCol(ProjectColumnHasProjectCheck, e.HasProjectCheck),
			handler.NewCol(ProjectColumnPrivateLabelingSetting, e.PrivateLabelingSetting),
			handler.NewCol(ProjectColumnState, domain.ProjectStateActive),
			handler.NewCol(ProjectColumnCreator, e.EditorUser()),
		},
	), nil
}

func (p *projectProjection) reduceProjectChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectChangeEvent)
	if !ok {
		logging.LogWithFields("HANDL-dk2iF", "seq", event.Sequence(), "expected", project.ProjectChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-s00Fs", "reduce.wrong.event.type")
	}
	if e.Name == nil && e.HasProjectCheck == nil && e.ProjectRoleAssertion == nil && e.ProjectRoleCheck == nil && e.PrivateLabelingSetting == nil {
		return crdb.NewNoOpStatement(e), nil
	}

	columns := make([]handler.Column, 0, 7)
	columns = append(columns, handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
		handler.NewCol(ProjectColumnSequence, e.Sequence()))
	if e.Name != nil {
		columns = append(columns, handler.NewCol(ProjectColumnName, *e.Name))
	}
	if e.ProjectRoleAssertion != nil {
		columns = append(columns, handler.NewCol(ProjectColumnProjectRoleAssertion, *e.ProjectRoleAssertion))
	}
	if e.ProjectRoleCheck != nil {
		columns = append(columns, handler.NewCol(ProjectColumnProjectRoleCheck, *e.ProjectRoleCheck))
	}
	if e.HasProjectCheck != nil {
		columns = append(columns, handler.NewCol(ProjectColumnHasProjectCheck, *e.HasProjectCheck))
	}
	if e.PrivateLabelingSetting != nil {
		columns = append(columns, handler.NewCol(ProjectColumnPrivateLabelingSetting, *e.PrivateLabelingSetting))
	}
	return crdb.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(ProjectColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *projectProjection) reduceProjectDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-8Nf2s", "seq", event.Sequence(), "expectedType", project.ProjectDeactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-LLp0f", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectColumnSequence, e.Sequence()),
			handler.NewCol(ProjectColumnState, domain.ProjectStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *projectProjection) reduceProjectReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-sm99f", "seq", event.Sequence(), "expectedType", project.ProjectReactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-9J98f", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectColumnSequence, e.Sequence()),
			handler.NewCol(ProjectColumnState, domain.ProjectStateActive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *projectProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-mL0sf", "seq", event.Sequence(), "expectedType", project.ProjectRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-5N9fs", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectColumnID, e.Aggregate().ID),
		},
	), nil
}
