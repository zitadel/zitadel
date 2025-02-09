package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/member"
	"github.com/zitadel/zitadel/internal/repository/user"

	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupMemberProjectionTable   = "projections.group_members14"
	GroupMemberGroupIDCol        = "group_id"
	GroupMemberUserIDCol         = "user_id"
	GroupMemberUserResourceOwner = "user_resource_owner"

	GroupMemberCreationDate  = "creation_date"
	GroupMemberChangeDate    = "change_date"
	GroupMemberSequence      = "sequence"
	GroupMemberResourceOwner = "resource_owner"
	GroupMemberInstanceID    = "instance_id"
)

var (
	groupMemberColumns = []*handler.InitColumn{
		handler.NewColumn(GroupMemberCreationDate, handler.ColumnTypeTimestamp),
		handler.NewColumn(GroupMemberChangeDate, handler.ColumnTypeTimestamp),
		handler.NewColumn(GroupMemberUserIDCol, handler.ColumnTypeText),
		handler.NewColumn(GroupMemberGroupIDCol, handler.ColumnTypeText),
		handler.NewColumn(GroupMemberUserResourceOwner, handler.ColumnTypeText),
		handler.NewColumn(GroupMemberSequence, handler.ColumnTypeInt64),
		handler.NewColumn(GroupMemberResourceOwner, handler.ColumnTypeText),
		handler.NewColumn(GroupMemberInstanceID, handler.ColumnTypeText),
	}
)

type reduceGroupMemberConfig struct {
	cols  []handler.Column
	conds []handler.Condition
}

type reduceGroupMemberOpt func(reduceGroupMemberConfig) reduceGroupMemberConfig

func withGroupMemberCol(col string, value interface{}) reduceGroupMemberOpt {
	return func(opt reduceGroupMemberConfig) reduceGroupMemberConfig {
		opt.cols = append(opt.cols, handler.NewCol(col, value))
		return opt
	}
}

func withGroupMemberCond(cond string, value interface{}) reduceGroupMemberOpt {
	return func(opt reduceGroupMemberConfig) reduceGroupMemberConfig {
		opt.conds = append(opt.conds, handler.NewCond(cond, value))
		return opt
	}
}

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
			groupMemberColumns, // handler.NewColumn(GroupMemberGroupIDCol, handler.ColumnTypeText),
			handler.NewPrimaryKey(GroupMemberInstanceID, GroupMemberGroupIDCol, GroupMemberUserIDCol),
			handler.WithIndex(handler.NewIndex("user_id", []string{GroupMemberUserIDCol})),
			handler.WithIndex(
				handler.NewIndex("gm_instance", []string{GroupMemberInstanceID},
					handler.WithInclude(
						GroupMemberCreationDate,
						GroupMemberChangeDate,
						GroupMemberSequence,
						GroupMemberResourceOwner,
					),
				),
			),
		),
	)
}

func reduceGroupMemberAdded(e member.MemberAddedEvent, userResourceOwner string, opts ...reduceMemberOpt) (*handler.Statement, error) {
	config := reduceMemberConfig{
		cols: []handler.Column{
			handler.NewCol(GroupMemberUserIDCol, e.UserID),
			// handler.NewCol(GroupMemberGroupIDCol, e.GroupID),
			handler.NewCol(GroupMemberUserResourceOwner, userResourceOwner),
			handler.NewCol(GroupMemberCreationDate, e.CreationDate()),
			handler.NewCol(GroupMemberChangeDate, e.CreationDate()),
			handler.NewCol(GroupMemberSequence, e.Sequence()),
			handler.NewCol(GroupMemberResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(GroupMemberInstanceID, e.Aggregate().InstanceID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.NewCreateStatement(&e, config.cols), nil
}

func reduceGroupMemberChanged(e member.MemberChangedEvent, opts ...reduceGroupMemberOpt) (*handler.Statement, error) {
	config := reduceGroupMemberConfig{
		cols: []handler.Column{
			handler.NewCol(GroupMemberChangeDate, e.CreationDate()),
			handler.NewCol(GroupMemberSequence, e.Sequence()),
		},
		conds: []handler.Condition{
			handler.NewCond(GroupMemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupMemberUserIDCol, e.UserID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.NewUpdateStatement(&e, config.cols, config.conds), nil
}

func reduceGroupMemberCascadeRemoved(e member.MemberCascadeRemovedEvent, opts ...reduceGroupMemberOpt) (*handler.Statement, error) {
	config := reduceGroupMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(GroupMemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupMemberUserIDCol, e.UserID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}
	return handler.NewDeleteStatement(&e, config.conds), nil
}

func reduceGroupMemberRemoved(e eventstore.Event, opts ...reduceGroupMemberOpt) (*handler.Statement, error) {
	config := reduceGroupMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(GroupMemberInstanceID, e.Aggregate().InstanceID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}
	return handler.NewDeleteStatement(e, config.conds), nil
}

func multiReduceGroupMemberOwnerRemoved(e eventstore.Event, opts ...reduceGroupMemberOpt) func(eventstore.Event) handler.Exec {
	config := reduceGroupMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(GroupMemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupMemberResourceOwner, e.Aggregate().ID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.AddDeleteStatement(
		config.conds,
	)
}

func groupMemberUserOwnerRemovedConds(e eventstore.Event, opts ...reduceGroupMemberOpt) []handler.Condition {
	config := reduceGroupMemberConfig{
		conds: []handler.Condition{
			handler.NewCond(GroupMemberInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupMemberUserResourceOwner, e.Aggregate().ID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}
	return config.conds
}

func reduceGroupMemberUserOwnerRemoved(e eventstore.Event, opts ...reduceGroupMemberOpt) (*handler.Statement, error) {
	return handler.NewDeleteStatement(
		e,
		groupMemberUserOwnerRemovedConds(e, opts...),
	), nil
}

func multiReduceGroupMemberUserOwnerRemoved(e eventstore.Event, opts ...reduceGroupMemberOpt) func(eventstore.Event) handler.Exec {
	return handler.AddDeleteStatement(
		groupMemberUserOwnerRemovedConds(e, opts...),
	)
}

func setGroupMemberContext(aggregate *eventstore.Aggregate) context.Context {
	return authz.WithInstanceID(context.Background(), aggregate.InstanceID)
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
	ctx := setGroupMemberContext(e.Aggregate())
	userOwner, err := getUserResourceOwner(ctx, g.es, e.Aggregate().InstanceID, e.UserID)
	if err != nil {
		return nil, err
	}
	return reduceGroupMemberAdded(
		*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID),
		userOwner,
		withMemberCol(GroupMemberGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MemberChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-10XK2", "reduce.wrong.event.type %s", group.MemberChangedType)
	}
	return reduceGroupMemberChanged(
		*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID),
		withGroupMemberCond(GroupMemberGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MemberCascadeRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-bHe54", "reduce.wrong.event.type %s", group.MemberCascadeRemovedType)
	}
	return reduceGroupMemberCascadeRemoved(
		*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
		withGroupMemberCond(GroupMemberGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.MemberRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-fKAOi", "reduce.wrong.event.type %s", group.MemberRemovedType)
	}
	return reduceGroupMemberRemoved(
		e,
		withGroupMemberCond(GroupMemberUserIDCol, e.UserID),
		withGroupMemberCond(GroupMemberGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-aYA60", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceGroupMemberRemoved(e,
		withGroupMemberCond(GroupMemberUserIDCol, e.Aggregate().ID),
	)
}

func (g *groupMemberProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-OHVFM", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewMultiStatement(
		e,
		multiReduceGroupMemberOwnerRemoved(e),
		multiReduceGroupMemberUserOwnerRemoved(e),
	), nil
}

func (g *groupMemberProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-OHVFM", "reduce.wrong.event.type %s", group.GroupRemovedType)
	}
	return reduceGroupMemberRemoved(e, withGroupMemberCond(GroupMemberGroupIDCol, e.Aggregate().ID))
}
