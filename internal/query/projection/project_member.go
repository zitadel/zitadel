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

type ProjectMemberProjection struct {
	crdb.StatementHandler
}

const (
	ProjectMemberProjectionTable = "zitadel.projections.project_members"
)

func NewProjectMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectMemberProjection {
	p := &ProjectMemberProjection{}
	config.ProjectionName = ProjectMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *ProjectMemberProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.MemberAddedType,
					Reduce: p.reduceOrgAdded,
				},
				{
					Event:  project.MemberChangedType,
					Reduce: p.reduceOrgChanged,
				},
				{
					Event:  project.MemberCascadeRemovedType,
					Reduce: p.reduceProjectMemberCascadeRemoved,
				},
				{
					Event:  project.MemberRemovedType,
					Reduce: p.reduceProjectMemberRemoved,
				},
			},
		},
	}
}

type ProjectMemberColumn string

const (
	ProjectMemberOrgIDCol  = "project_id"
	ProjectMemberUserIDCol = "user_id"
	ProjectMemberRolesCol  = "roles"
)

func (p *ProjectMemberProjection) reduceOrgAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-vHSJA", "seq", event.Sequence(), "expectedType", project.MemberAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-7Ybl2", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectMemberUserIDCol, e.UserID),
			handler.NewCol(ProjectMemberRolesCol, e.Roles),
		},
	), nil
}

func (p *ProjectMemberProjection) reduceOrgChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-RVPu0", "seq", event.Sequence(), "expected", project.MemberChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-O7hHz", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectMemberRolesCol, e.Roles),
		},
		[]handler.Condition{
			handler.NewCond(ProjectMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(ProjectMemberUserIDCol, e.UserID),
		},
	), nil
}

func (p *ProjectMemberProjection) reduceProjectMemberCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-kEKtY", "seq", event.Sequence(), "expected", project.MemberCascadeRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-lx4PS", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(ProjectMemberUserIDCol, e.UserID),
		},
	), nil
}

func (p *ProjectMemberProjection) reduceProjectMemberRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-O03Cf", "seq", event.Sequence(), "expected", project.MemberRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-IZBJc", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectMemberOrgIDCol, e.Aggregate().ResourceOwner),
			handler.NewCond(ProjectMemberUserIDCol, e.UserID),
		},
	), nil
}
