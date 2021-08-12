package projection

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/project"
)

type ProjectProjection struct {
	crdb.StatementHandler
}

func NewProjectProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectProjection {
	p := &ProjectProjection{}
	config.ProjectionName = "projections.projects"
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

type projectState int8

const (
	projectIDCol           = "id"
	projectNameCol         = "name"
	projectCreationDateCol = "creation_date"
	projectChangeDateCol   = "change_date"
	projectOwnerCol        = "owner_id"
	projectCreatorCol      = "creator_id"
	projectStateCol        = "state"

	projectActive projectState = iota
	projectInactive
)

func (p *ProjectProjection) reduceProjectAdded(event eventstore.EventReader) ([]handler.Statement, error) {
	e := event.(*project.ProjectAddedEvent)

	return []handler.Statement{
		crdb.NewCreateStatement(
			e,
			[]handler.Column{
				handler.NewCol(projectIDCol, e.Aggregate().ID),
				handler.NewCol(projectNameCol, e.Name),
				handler.NewCol(projectCreationDateCol, e.CreationDate()),
				handler.NewCol(projectChangeDateCol, e.CreationDate()),
				handler.NewCol(projectOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(projectCreatorCol, e.EditorUser()),
				handler.NewCol(projectStateCol, projectActive),
			},
		),
	}, nil
}

func (p *ProjectProjection) reduceProjectChanged(event eventstore.EventReader) ([]handler.Statement, error) {
	e := event.(*project.ProjectChangeEvent)

	if e.Name == nil {
		return []handler.Statement{crdb.NewNoOpStatement(e)}, nil
	}

	return []handler.Statement{
		crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(projectNameCol, e.Name),
				handler.NewCol(projectChangeDateCol, e.CreationDate()),
			},
			[]handler.Column{
				handler.NewCol(projectIDCol, e.Aggregate().ID),
			},
		),
	}, nil
}

func (p *ProjectProjection) reduceProjectDeactivated(event eventstore.EventReader) ([]handler.Statement, error) {
	e := event.(*project.ProjectDeactivatedEvent)

	return []handler.Statement{
		crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(projectStateCol, projectInactive),
				handler.NewCol(projectChangeDateCol, e.CreationDate()),
			},
			[]handler.Column{
				handler.NewCol(projectIDCol, e.Aggregate().ID),
			},
		),
	}, nil
}

func (p *ProjectProjection) reduceProjectReactivated(event eventstore.EventReader) ([]handler.Statement, error) {
	e := event.(*project.ProjectReactivatedEvent)

	return []handler.Statement{
		crdb.NewUpdateStatement(
			e,
			[]handler.Column{
				handler.NewCol(projectStateCol, projectActive),
				handler.NewCol(projectChangeDateCol, e.CreationDate()),
			},
			[]handler.Column{
				handler.NewCol(projectIDCol, e.Aggregate().ID),
			},
		),
	}, nil
}

func (p *ProjectProjection) reduceProjectRemoved(event eventstore.EventReader) ([]handler.Statement, error) {
	e := event.(*project.ProjectRemovedEvent)

	return []handler.Statement{
		crdb.NewDeleteStatement(
			e,
			[]handler.Column{
				handler.NewCol(projectIDCol, e.Aggregate().ID),
			},
		),
	}, nil
}
