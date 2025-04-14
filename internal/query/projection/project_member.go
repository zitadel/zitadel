package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ProjectMemberProjectionTable = "projections.project_members4"
	ProjectMemberProjectIDCol    = "project_id"
)

type projectMemberProjection struct {
	es handler.EventStore
}

func newProjectMemberProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &projectMemberProjection{es: config.Eventstore})
}

func (*projectMemberProjection) Name() string {
	return ProjectMemberProjectionTable
}

func (*projectMemberProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable(
			append(memberColumns,
				handler.NewColumn(ProjectMemberProjectIDCol, handler.ColumnTypeText),
			),
			handler.NewPrimaryKey(MemberInstanceID, ProjectMemberProjectIDCol, MemberUserIDCol),
			handler.WithIndex(handler.NewIndex("user_id", []string{MemberUserIDCol})),
			handler.WithIndex(
				handler.NewIndex("pm_instance", []string{MemberInstanceID},
					handler.WithInclude(
						MemberCreationDate,
						MemberChangeDate,
						MemberRolesCol,
						MemberSequence,
						MemberResourceOwner,
					),
				),
			),
		),
	)
}

func (p *projectMemberProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.MemberAddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  project.MemberChangedEventType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  project.MemberCascadeRemovedEventType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  project.MemberRemovedEventType,
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
			EventReducers: []handler.EventReducer{
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOrgRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-bgx5Q", "reduce.wrong.event.type %s", project.MemberAddedEventType)
	}
	ctx := setMemberContext(e.Aggregate())
	userOwner, err := getUserResourceOwner(ctx, p.es, e.Aggregate().InstanceID, e.UserID)
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-90WJ1", "reduce.wrong.event.type %s", project.MemberChangedEventType)
	}
	return reduceMemberChanged(
		*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *projectMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberCascadeRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-aGd43", "reduce.wrong.event.type %s", project.MemberCascadeRemovedEventType)
	}
	return reduceMemberCascadeRemoved(
		*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *projectMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MemberRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-eJZPh", "reduce.wrong.event.type %s", project.MemberRemovedEventType)
	}
	return reduceMemberRemoved(e,
		withMemberCond(MemberUserIDCol, e.UserID),
		withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID),
	)
}

func (p *projectMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-aYA60", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *projectMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-NGUEL", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewMultiStatement(
		e,
		multiReduceMemberOwnerRemoved(e),
		multiReduceMemberUserOwnerRemoved(e),
	), nil
}

func (p *projectMemberProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-NGUEL", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(ProjectMemberProjectIDCol, e.Aggregate().ID))
}
