package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupMemberProjectionTable = "projections.group_members"
	GroupMemberGroupIDCol      = "group_id"
)

type groupMemberProjection struct {
	es handler.EventStore
}

func newGroupMemberProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &groupMemberProjection{es: config.Eventstore})
}

func (*groupMemberProjection) Name() string {
	return GroupMemberProjectionTable
}

func (*groupMemberProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable(
			append(memberColumns,
				handler.NewColumn(GroupMemberGroupIDCol, handler.ColumnTypeText),
			),
			handler.NewPrimaryKey(MemberInstanceID, GroupMemberGroupIDCol, MemberUserIDCol),
			handler.WithIndex(handler.NewIndex("user_id", []string{MemberUserIDCol})),
			handler.WithIndex(
				handler.NewIndex("gm_instance", []string{MemberInstanceID},
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

func (g *groupMemberProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.MemberAddedType,
					Reduce: g.reduceAdded,
				},
				{
					Event:  group.MemberChangedType,
					Reduce: g.reduceChanged,
				},
				{
					Event:  group.MemberCascadeRemovedType,
					Reduce: g.reduceCascadeRemoved,
				},
				{
					Event:  group.MemberRemovedType,
					Reduce: g.reduceRemoved,
				},
				{
					Event:  group.GroupRemovedType,
					Reduce: g.reduceGroupRemoved,
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.UserRemovedType,
					Reduce: g.reduceUserRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: g.reduceOrgRemoved,
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

func (g *groupMemberProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MemberAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-chy6O", "reduce.wrong.event.type %s", group.MemberAddedType)
	}
	ctx := setMemberContext(e.Aggregate())
	userOwner, err := getUserResourceOwner(ctx, g.es, e.Aggregate().InstanceID, e.UserID)
	if err != nil {
		return nil, err
	}
	return reduceMemberAdded(
		*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		userOwner,
		withMemberCol(GroupMemberGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MemberChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-10XK2", "reduce.wrong.event.type %s", group.MemberChangedType)
	}
	return reduceMemberChanged(
		*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCond(GroupMemberGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MemberCascadeRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-bHe54", "reduce.wrong.event.type %s", group.MemberCascadeRemovedType)
	}
	return reduceMemberCascadeRemoved(
		*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(GroupMemberGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MemberRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-fKAOi", "reduce.wrong.event.type %s", group.MemberRemovedType)
	}
	return reduceMemberRemoved(e,
		withMemberCond(MemberUserIDCol, e.UserID),
		withMemberCond(GroupMemberGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-bZB60", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(MemberUserIDCol, e.Aggregate().ID))
}

func (g *groupMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-OHVFM", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewMultiStatement(
		e,
		multiReduceMemberOwnerRemoved(e),
		multiReduceMemberUserOwnerRemoved(e),
	), nil
}

func (g *groupMemberProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-OHVFM", "reduce.wrong.event.type %s", group.GroupRemovedType)
	}
	return reduceMemberRemoved(e, withMemberCond(GroupMemberGroupIDCol, e.Aggregate().ID))
}
