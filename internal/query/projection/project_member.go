package projection

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/member"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/project"
	"github.com/caos/zitadel/internal/repository/user"
)

const (
	ProjectMemberProjectionTable = "projections.project_members"
	ProjectMemberProjectIDCol    = "project_id"
)

type ProjectMemberProjection struct {
	crdb.StatementHandler
}

func NewProjectMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectMemberProjection {
	p := new(ProjectMemberProjection)
	config.ProjectionName = ProjectMemberProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable(
			append(memberColumns,
				crdb.NewColumn(ProjectMemberProjectIDCol, crdb.ColumnTypeText),
			),
			crdb.NewPrimaryKey(ProjectMemberProjectIDCol, MemberUserIDCol),
			crdb.NewIndex("user_idx", []string{MemberUserIDCol}),
		),
	)

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

func (p *ProjectMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-bgx5Q", "reduce.wrong.event.type %s", project.MemberAddedType)
	}
	return reduceMemberAdded(
		*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCol(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *ProjectMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-90WJ1", "reduce.wrong.event.type %s", project.MemberChangedType)
	}
	return reduceMemberChanged(
		*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *ProjectMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberCascadeRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-aGd43", "reduce.wrong.event.type %s", project.MemberCascadeRemovedType)
	}
	return reduceMemberCascadeRemoved(
		*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *ProjectMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-eJZPh", "reduce.wrong.event.type %s", project.MemberRemovedType)
	}
	return reduceMemberRemoved(e,
		withMemberCond(MemberUserIDCol, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *ProjectMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-aYA60", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *ProjectMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	//TODO: as soon as org deletion is implemented:
	// Case: The user has resource owner A and project has resource owner B
	// if org B deleted it works
	// if org A is deleted, the membership wouldn't be deleted
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-NGUEL", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberResourceOwner, e.Aggregate().ID))
}

func (p *ProjectMemberProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-NGUEL", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID))
}
