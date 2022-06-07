package projection

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type projectMemberProjection struct {
	crdb.StatementHandler
}

const (
	ProjectMemberProjectionTable = "zitadel.projections.project_members"
)

func newProjectMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *projectMemberProjection {
	p := &projectMemberProjection{}
	config.ProjectionName = ProjectMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *projectMemberProjection) reducers() []handler.AggregateReducer {
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
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
				},
			},
		},
	}
}

type ProjectMemberColumn string

const (
	ProjectMemberProjectIDCol = "project_id"
)

func (p *projectMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
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

func (p *projectMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *projectMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
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

func (p *projectMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-X0yvM", "seq", event.Sequence(), "expectedType", project.MemberRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-eJZPh", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e,
		withMemberCond(MemberUserIDCol, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *projectMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-g8eWd", "seq", event.Sequence(), "expected", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-aYA60", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *projectMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	//TODO: as soon as org deletion is implemented:
	// Case: The user has resource owner A and project has resource owner B
	// if org B deleted it works
	// if org A is deleted, the membership wouldn't be deleted
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-q7H8D", "seq", event.Sequence(), "expected", org.OrgRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-NGUEL", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberResourceOwner, e.Aggregate().ID))
}

func (p *projectMemberProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-q7H8D", "seq", event.Sequence(), "expected", project.ProjectRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-NGUEL", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID))
}
