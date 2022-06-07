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

type projectGrantMemberProjection struct {
	crdb.StatementHandler
}

const (
	ProjectGrantMemberProjectionTable = "zitadel.projections.project_grant_members"
)

func newProjectGrantMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *projectGrantMemberProjection {
	p := &projectGrantMemberProjection{}
	config.ProjectionName = ProjectGrantMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *projectGrantMemberProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.GrantMemberAddedType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  project.GrantMemberChangedType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  project.GrantMemberCascadeRemovedType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  project.GrantMemberRemovedType,
					Reduce: p.reduceRemoved,
				},
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
				{
					Event:  project.GrantRemovedType,
					Reduce: p.reduceProjectGrantRemoved,
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

type ProjectGrantMemberColumn string

const (
	ProjectGrantMemberProjectIDCol = "project_id"
	ProjectGrantMemberGrantIDCol   = "grant_id"
)

func (p *projectGrantMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-csr8B", "seq", event.Sequence(), "expectedType", project.GrantMemberAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-0EBQf", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(
		*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCol(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCol(ProjectGrantMemberGrantIDCol, e.GrantID),
	)
}

func (p *projectGrantMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-ZubbI", "seq", event.Sequence(), "expectedType", project.GrantMemberChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-YX5Tk", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(
		*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCond(ProjectGrantMemberGrantIDCol, e.GrantID),
	)
}

func (p *projectGrantMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-azx7K", "seq", event.Sequence(), "expectedType", project.GrantMemberCascadeRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-adnHG", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(
		*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCond(ProjectGrantMemberGrantIDCol, e.GrantID),
	)
}

func (p *projectGrantMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-6Z4dH", "seq", event.Sequence(), "expectedType", project.GrantMemberRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-MGNnA", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e,
		withMemberCond(MemberUserIDCol, e.UserID),
		withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCond(ProjectGrantMemberGrantIDCol, e.GrantID),
	)
}

func (p *projectGrantMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-UVMmT", "seq", event.Sequence(), "expected", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-rufJr", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *projectGrantMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	//TODO: as soon as org deletion is implemented:
	// Case: The user has resource owner A and project has resource owner B
	// if org B deleted it works
	// if org A is deleted, the membership wouldn't be deleted
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Sq9FV", "seq", event.Sequence(), "expected", org.OrgRemovedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Zzp6o", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(MemberResourceOwner, e.Aggregate().ID))
}

func (p *projectGrantMemberProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-sGmCA", "seq", event.Sequence(), "expected", project.ProjectRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-JLODy", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e, withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID))
}

func (p *projectGrantMemberProjection) reduceProjectGrantRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-sHabO", "seq", event.Sequence(), "expected", project.GrantRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-D1J9R", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(e,
		withMemberCond(ProjectGrantMemberGrantIDCol, e.GrantID),
		withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
	)
}
