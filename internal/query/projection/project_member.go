package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	ProjectMemberProjectionTable = "projections.project_members3"
	ProjectMemberProjectIDCol    = "project_id"
)

type projectMemberProjection struct {
	crdb.StatementHandler
}

func newProjectMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *projectMemberProjection {
	p := new(projectMemberProjection)
	config.ProjectionName = ProjectMemberProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable(
			append(memberColumns,
				crdb.NewColumn(ProjectMemberProjectIDCol, crdb.ColumnTypeText),
			),
			crdb.NewPrimaryKey(MemberInstanceID, ProjectMemberProjectIDCol, MemberUserIDCol),
			crdb.WithIndex(crdb.NewIndex("user_id", []string{MemberUserIDCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{MemberOwnerRemoved})),
			crdb.WithIndex(crdb.NewIndex("user_owner_removed", []string{MemberUserOwnerRemoved})),
		),
	)

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
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(MemberInstanceID),
				},
			},
		},
	}
}

func (p *projectMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-bgx5Q", "reduce.wrong.event.type %s", project.MemberAddedType)
	}
	ctx := setMemberContext(e.Aggregate())
	userOwner, err := getResourceOwnerOfUser(ctx, p.Eventstore, e.Aggregate().InstanceID, e.UserID)
	if err != nil {
		return nil, err
	}
	return reduceMemberAdded(
		*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		userOwner,
		withMemberCol(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *projectMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-90WJ1", "reduce.wrong.event.type %s", project.MemberChangedType)
	}
	return reduceMemberChanged(
		*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *projectMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberCascadeRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-aGd43", "reduce.wrong.event.type %s", project.MemberCascadeRemovedType)
	}
	return reduceMemberCascadeRemoved(
		*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *projectMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-eJZPh", "reduce.wrong.event.type %s", project.MemberRemovedType)
	}
	return reduceMemberRemoved(e,
		withMemberCond(MemberUserIDCol, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *projectMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-aYA60", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *projectMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-NGUEL", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return crdb.NewMultiStatement(
		e,
		multiReduceMemberOwnerRemoved(e),
		multiReduceMemberUserOwnerRemoved(e),
	), nil
}

func (p *projectMemberProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-NGUEL", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID))
}
