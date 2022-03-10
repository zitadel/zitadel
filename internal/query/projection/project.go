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

const (
	ProjectProjectionTable = "zitadel.projections.projects"

	ProjectColumnID                     = "id"
	ProjectColumnCreationDate           = "creation_date"
	ProjectColumnChangeDate             = "change_date"
	ProjectColumnSequence               = "sequence"
	ProjectColumnState                  = "state"
	ProjectColumnResourceOwner          = "resource_owner"
	ProjectColumnName                   = "name"
	ProjectColumnProjectRoleAssertion   = "project_role_assertion"
	ProjectColumnProjectRoleCheck       = "project_role_check"
	ProjectColumnHasProjectCheck        = "has_project_check"
	ProjectColumnPrivateLabelingSetting = "private_labeling_setting"
	ProjectColumnCreator                = "creator_id" //TODO: necessary?
)

type ProjectProjection struct {
	crdb.StatementHandler
}

func NewProjectProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectProjection {
	p := new(ProjectProjection)
	config.ProjectionName = ProjectProjectionTable
	config.Reducers = p.reducers()
	config.InitChecks = []*handler.Check{
		crdb.NewTableCheck(
			crdb.NewTable([]*crdb.Column{
				crdb.NewColumn(ProjectColumnID, crdb.ColumnTypeText),
				crdb.NewColumn(ProjectColumnCreationDate, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(ProjectColumnChangeDate, crdb.ColumnTypeTimestamp),
				crdb.NewColumn(ProjectColumnSequence, crdb.ColumnTypeInt64),
				crdb.NewColumn(ProjectColumnState, crdb.ColumnTypeEnum),
				crdb.NewColumn(ProjectColumnResourceOwner, crdb.ColumnTypeText),
				crdb.NewColumn(ProjectColumnName, crdb.ColumnTypeText),
				crdb.NewColumn(ProjectColumnProjectRoleAssertion, crdb.ColumnTypeBool),
				crdb.NewColumn(ProjectColumnProjectRoleCheck, crdb.ColumnTypeBool),
				crdb.NewColumn(ProjectColumnHasProjectCheck, crdb.ColumnTypeBool),
				crdb.NewColumn(ProjectColumnPrivateLabelingSetting, crdb.ColumnTypeEnum),
				crdb.NewColumn(ProjectColumnCreator, crdb.ColumnTypeText),
			},
				crdb.NewPrimaryKey(ProjectColumnID),
			),
		),
	}
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

func (p *ProjectProjection) reduceProjectAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *ProjectProjection) reduceProjectChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *ProjectProjection) reduceProjectDeactivated(event eventstore.Event) (*handler.Statement, error) {
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

func (p *ProjectProjection) reduceProjectReactivated(event eventstore.Event) (*handler.Statement, error) {
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

func (p *ProjectProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
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
