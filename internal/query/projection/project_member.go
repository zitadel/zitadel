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
					Event:  project.GrantMemberAddedType,
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
	ProjectMemberGrantIDCol   = "grant_id"
)

func (p *ProjectMemberProjection) reduceAdded(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *project.MemberAddedEvent:
		return reduceMemberAdded(e.MemberAddedEvent, withMemberCol(ProjectMemberProjectIDCol, e.Aggregate().ID))
	case *project.GrantMemberAddedEvent:
		return reduceMemberAdded(
			*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID, e.Roles...),
			withMemberCol(ProjectMemberProjectIDCol, e.Aggregate().ID),
			withMemberCol(ProjectMemberGrantIDCol, e.GrantID),
		)
	default:
		logging.LogWithFields("HANDL-tPdUI", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{project.MemberAddedType, project.GrantMemberAddedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-rddM1", "reduce.wrong.event.type")
	}
}

func (p *ProjectMemberProjection) reduceChanged(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *project.MemberChangedEvent:
		return reduceMemberChanged(e.MemberChangedEvent, withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID))
	case *project.GrantMemberChangedEvent:
		return reduceMemberChanged(
			*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID, e.Roles...),
			withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
			withMemberCond(ProjectMemberGrantIDCol, e.GrantID),
		)
	default:
		logging.LogWithFields("HANDL-LxWSn", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{project.MemberChangedType, project.GrantMemberChangedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-FM8L7", "reduce.wrong.event.type")
	}
}

func (p *ProjectMemberProjection) reduceCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *project.MemberCascadeRemovedEvent:
		return reduceMemberCascadeRemoved(e.MemberCascadeRemovedEvent, withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID))
	case *project.GrantMemberCascadeRemovedEvent:
		return reduceMemberCascadeRemoved(
			*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
			withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
			withMemberCond(ProjectMemberGrantIDCol, e.GrantID),
		)
	default:
		logging.LogWithFields("HANDL-6gFXG", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{project.MemberCascadeRemovedType, project.GrantMemberCascadeRemovedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-lwmfd", "reduce.wrong.event.type")
	}
}

func (p *ProjectMemberProjection) reduceRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	switch e := event.(type) {
	case *project.MemberRemovedEvent:
		return reduceMemberRemoved(e.MemberRemovedEvent, withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID))
	case *project.GrantMemberRemovedEvent:
		return reduceMemberRemoved(
			*member.NewRemovedEvent(&e.BaseEvent, e.UserID),
			withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
			withMemberCond(ProjectMemberGrantIDCol, e.GrantID),
		)
	default:
		logging.LogWithFields("HANDL-ryUlI", "seq", event.Sequence(), "expectedTypes", []eventstore.EventType{project.MemberRemovedType, project.GrantMemberRemovedType}).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-vU5sx", "reduce.wrong.event.type")
	}
}
