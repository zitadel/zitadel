package projection

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/project"
)

type ProjectGrantProjection struct {
	crdb.StatementHandler
}

func NewProjectGrantProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectGrantProjection {
	p := &ProjectGrantProjection{}
	config.ProjectionName = "zitadel.projections.project_grants"
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *ProjectGrantProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.GrantAddedType,
					Reduce: p.reduceProjectGrantAdded,
				},
				{
					Event:  project.GrantChangedType,
					Reduce: p.reduceProjectGrantChanged,
				},
				{
					Event:  project.GrantCascadeChangedType,
					Reduce: p.reduceProjectGrantCascadeChanged,
				},
				{
					Event:  project.GrantDeactivatedType,
					Reduce: p.reduceProjectGrantDeactivated,
				},
				{
					Event:  project.GrantReactivatedType,
					Reduce: p.reduceProjectGrantReactivated,
				},
				{
					Event:  project.GrantRemovedType,
					Reduce: p.reduceProjectGrantRemoved,
				},
			},
		},
	}
}

type projectGrantState int8

const (
	projectGrantIDCol           = "grant_id"
	projectGrantProjectIDCol    = "project_id"
	projectGrantGrantedOrgIDCol = "granted_org_id"
	projectGrantRoleKeysCol     = "role_keys"
	projectGrantCreationDateCol = "creation_date"
	projectGrantChangeDateCol   = "change_date"
	projectGrantOwnerCol        = "owner_id"
	projectGrantCreatorCol      = "creator_id"
	projectGrantStateCol        = "state"

	projectGrantActive projectGrantState = iota
	projectGrantInactive
)

func (p *ProjectGrantProjection) reduceProjectGrantAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e := event.(*project.GrantAddedEvent)

	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectGrantProjectIDCol, e.Aggregate().ID),
			handler.NewCol(projectGrantIDCol, e.GrantID),
			handler.NewCol(projectGrantGrantedOrgIDCol, e.GrantedOrgID),
			handler.NewCol(projectGrantRoleKeysCol, e.RoleKeys),
			handler.NewCol(projectGrantCreationDateCol, e.CreationDate()),
			handler.NewCol(projectGrantChangeDateCol, e.CreationDate()),
			handler.NewCol(projectGrantOwnerCol, e.Aggregate().ResourceOwner),
			handler.NewCol(projectGrantCreatorCol, e.EditorUser()),
			handler.NewCol(projectGrantStateCol, projectGrantActive),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e := event.(*project.GrantChangedEvent)

	if e.RoleKeys == nil {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectGrantRoleKeysCol, e.RoleKeys),
			handler.NewCol(projectChangeDateCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(projectGrantIDCol, e.GrantID),
			handler.NewCond(projectGrantProjectIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantCascadeChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e := event.(*project.GrantChangedEvent)

	if e.RoleKeys == nil {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectGrantRoleKeysCol, e.RoleKeys),
			handler.NewCol(projectGrantChangeDateCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(projectGrantIDCol, e.GrantID),
			handler.NewCond(projectGrantProjectIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantDeactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e := event.(*project.GrantDeactivateEvent)

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectGrantStateCol, projectGrantInactive),
			handler.NewCol(projectGrantChangeDateCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(projectGrantIDCol, e.GrantID),
			handler.NewCond(projectGrantProjectIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantReactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e := event.(*project.GrantReactivatedEvent)

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(projectGrantStateCol, projectGrantActive),
			handler.NewCol(projectGrantChangeDateCol, e.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(projectGrantIDCol, e.GrantID),
			handler.NewCond(projectGrantProjectIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e := event.(*project.GrantRemovedEvent)

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(projectGrantIDCol, e.GrantID),
			handler.NewCond(projectGrantProjectIDCol, e.Aggregate().ID),
		},
	), nil
}
