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
	ProjectGrantMemberProjectionTable = "projections.project_grant_members4"
	ProjectGrantMemberProjectIDCol    = "project_id"
	ProjectGrantMemberGrantIDCol      = "grant_id"
	ProjectGrantMemberGrantedOrg      = "granted_org"
)

type projectGrantMemberProjection struct {
	es handler.EventStore
}

func newProjectGrantMemberProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &projectGrantMemberProjection{es: config.Eventstore})
}

func (*projectGrantMemberProjection) Name() string {
	return ProjectGrantMemberProjectionTable
}

func (*projectGrantMemberProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable(
			append(memberColumns,
				handler.NewColumn(ProjectGrantMemberProjectIDCol, handler.ColumnTypeText),
				handler.NewColumn(ProjectGrantMemberGrantIDCol, handler.ColumnTypeText),
				handler.NewColumn(ProjectGrantMemberGrantedOrg, handler.ColumnTypeText),
			),
			handler.NewPrimaryKey(MemberInstanceID, ProjectGrantMemberProjectIDCol, ProjectGrantMemberGrantIDCol, MemberUserIDCol),
			handler.WithIndex(handler.NewIndex("user_id", []string{MemberUserIDCol})),
			handler.WithIndex(
				handler.NewIndex("pgm_instance", []string{MemberInstanceID},
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

func (p *projectGrantMemberProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
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

func (p *projectGrantMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-0EBQf", "reduce.wrong.event.type %s", project.GrantMemberAddedType)
	}
	ctx := setMemberContext(e.Aggregate())
	userOwner, err := getResourceOwnerOfUser(ctx, p.es, e.Aggregate().InstanceID, e.UserID)
	if err != nil {
		return nil, err
	}
	grantedOrg, err := getGrantedOrgOfGrantedProject(ctx, p.es, e.Aggregate().InstanceID, e.Aggregate().ID, e.GrantID)
	if err != nil {
		return nil, err
	}
	return reduceMemberAdded(
		*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		userOwner,
		withMemberCol(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCol(ProjectGrantMemberGrantIDCol, e.GrantID),
		withMemberCol(ProjectGrantMemberGrantedOrg, grantedOrg),
	)
}

func (p *projectGrantMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-YX5Tk", "reduce.wrong.event.type %s", project.GrantMemberChangedType)
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-adnHG", "reduce.wrong.event.type %s", project.GrantMemberCascadeRemovedType)
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-MGNnA", "reduce.wrong.event.type %s", project.GrantMemberRemovedType)
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-rufJr", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (p *projectGrantMemberProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Z2p6o", "reduce.wrong.event.type %s", instance.InstanceRemovedEventType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberInstanceID, e.Aggregate().ID))
}

func (p *projectGrantMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Zzp6o", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewMultiStatement(
		e,
		multiReduceMemberOwnerRemoved(e),
		multiReduceMemberUserOwnerRemoved(e),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(ProjectGrantMemberGrantedOrg, e.Aggregate().ID),
			},
		),
	), nil
}

func (p *projectGrantMemberProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JLODy", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID))
}

func (p *projectGrantMemberProjection) reduceProjectGrantRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-D1J9R", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}
	return reduceMemberRemoved(e,
		withMemberCond(ProjectGrantMemberGrantIDCol, e.GrantID),
		withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
	)
}
