package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/member"
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
		logging.LogWithFields("HANDL-3FRys", "seq", event.Sequence(), "expectedType", project.MemberAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-bgx5Q", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(
		*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCol(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *ProjectMemberProjection) reduceChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-9hgMR", "seq", event.Sequence(), "expectedType", project.MemberChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-90WJ1", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(
		*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *ProjectMemberProjection) reduceCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-2kyYa", "seq", event.Sequence(), "expectedType", project.MemberCascadeRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-aGd43", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(
		*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *ProjectMemberProjection) reduceRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-X0yvM", "seq", event.Sequence(), "expectedType", project.MemberRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-eJZPh", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(
		*member.NewRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}
