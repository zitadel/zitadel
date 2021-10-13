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
	"github.com/lib/pq"
)

type ProjectGrantProjection struct {
	crdb.StatementHandler
}

const ProjectGrantProjectionTable = "zitadel.projections.project_grants"

func NewProjectGrantProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectGrantProjection {
	p := &ProjectGrantProjection{}
	config.ProjectionName = ProjectGrantProjectionTable
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
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
	}
}

const (
	ProjectGrantColumnProjectID     = "project_id"
	ProjectGrantColumnGrantID       = "grant_id"
	ProjectGrantColumnCreationDate  = "creation_date"
	ProjectGrantColumnChangeDate    = "change_date"
	ProjectGrantColumnResourceOwner = "resource_owner"
	ProjectGrantColumnState         = "state"
	ProjectGrantColumnSequence      = "sequence"
	ProjectGrantColumnGrantedOrgID  = "granted_org_id"
	ProjectGrantColumnRoleKeys      = "granted_role_keys"
	ProjectGrantColumnCreator       = "creator_id"
)

func (p *ProjectGrantProjection) reduceProjectGrantAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Mi4g9", "seq", event.Sequence(), "expectedType", project.GrantAddedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-g92Fg", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCol(ProjectGrantColumnProjectID, e.Aggregate().ID),
			handler.NewCol(ProjectGrantColumnCreationDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectGrantColumnState, domain.ProjectGrantStateActive),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnGrantedOrgID, e.GrantedOrgID),
			handler.NewCol(ProjectGrantColumnRoleKeys, pq.StringArray(e.RoleKeys)),
			handler.NewCol(ProjectGrantColumnCreator, e.EditorUser()),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-M00fH", "seq", event.Sequence(), "expectedType", project.GrantChangedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-g0fg4", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnRoleKeys, pq.StringArray(e.RoleKeys)),
		},
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantCascadeChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-K0fwR", "seq", event.Sequence(), "expectedType", project.GrantCascadeChangedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-ll9Ts", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnRoleKeys, e.RoleKeys),
		},
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantDeactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantDeactivateEvent)
	if !ok {
		logging.LogWithFields("HANDL-Ple9f", "seq", event.Sequence(), "expectedType", project.GrantDeactivatedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-0fj2f", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnState, domain.ProjectGrantStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantReactivated(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Ip0hr", "seq", event.Sequence(), "expectedType", project.GrantReactivatedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-2M0ve", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnState, domain.ProjectGrantStateActive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectGrantRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-M0pfs", "seq", event.Sequence(), "expectedType", project.GrantRemovedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-o0w4f", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
		},
	), nil
}

func (p *ProjectGrantProjection) reduceProjectRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Ms0fe", "seq", event.Sequence(), "expectedType", project.ProjectRemovedType).Error("was not an  event")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-gn9rw", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
		},
	), nil
}
