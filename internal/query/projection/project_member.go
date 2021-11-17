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
					Reduce: p.reduceAdded,
				},
				{
					Event:  project.MemberChangedType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  project.MemberCascadeRemovedType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  project.MemberRemovedType,
					Reduce: p.reduceRemoved,
				},
			},
		},
	}
}

type ProjectMemberColumn string

const (
	ProjectMemberProjectIDCol = "project_id"
)

func (p *ProjectMemberProjection) reduceAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-CEmLi", "seq", event.Sequence(), "expectedType", project.MemberAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-u92W1", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(e.MemberAddedEvent, ProjectMemberProjectIDCol)
}

func (p *ProjectMemberProjection) reduceChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-2mexJ", "seq", event.Sequence(), "expected", project.MemberChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-fl6sR", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(e.MemberChangedEvent, ProjectMemberProjectIDCol)
}

func (p *ProjectMemberProjection) reduceCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Madqe", "seq", event.Sequence(), "expected", project.MemberCascadeRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-byUCI", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent, ProjectMemberProjectIDCol)
}

func (p *ProjectMemberProjection) reduceRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-MBL7R", "seq", event.Sequence(), "expected", project.MemberRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-lnEzp", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e.MemberRemovedEvent, ProjectMemberProjectIDCol)
}
