package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	groupusers "github.com/zitadel/zitadel/internal/repository/group_users"
	"github.com/zitadel/zitadel/internal/repository/user"

	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupUserProjectionTable   = "projections.group_users"
	GroupUserGroupIDCol        = "group_id"
	GroupUserUserIDCol         = "user_id"
	GroupUserUserResourceOwner = "user_resource_owner"

	GroupUserCreationDate  = "creation_date"
	GroupUserChangeDate    = "change_date"
	GroupUserSequence      = "sequence"
	GroupUserResourceOwner = "resource_owner"
	GroupUserInstanceID    = "instance_id"
	GroupUserAttributes    = "attributes"
)

var (
	groupUserColumns = []*handler.InitColumn{
		handler.NewColumn(GroupUserCreationDate, handler.ColumnTypeTimestamp),
		handler.NewColumn(GroupUserChangeDate, handler.ColumnTypeTimestamp),
		handler.NewColumn(GroupUserUserIDCol, handler.ColumnTypeText),
		handler.NewColumn(GroupUserGroupIDCol, handler.ColumnTypeText),
		handler.NewColumn(GroupUserUserResourceOwner, handler.ColumnTypeText),
		handler.NewColumn(GroupUserSequence, handler.ColumnTypeInt64),
		handler.NewColumn(GroupUserResourceOwner, handler.ColumnTypeText),
		handler.NewColumn(GroupUserInstanceID, handler.ColumnTypeText),
		handler.NewColumn(GroupUserAttributes, handler.ColumnTypeTextArray, handler.Nullable()),
	}
)

type reduceGroupUserConfig struct {
	cols  []handler.Column
	conds []handler.Condition
}

type reduceGroupUserOpt func(reduceGroupUserConfig) reduceGroupUserConfig

func withGroupUserCol(col string, value interface{}) reduceGroupUserOpt {
	return func(opt reduceGroupUserConfig) reduceGroupUserConfig {
		opt.cols = append(opt.cols, handler.NewCol(col, value))
		return opt
	}
}

func withGroupUserCond(cond string, value interface{}) reduceGroupUserOpt {
	return func(opt reduceGroupUserConfig) reduceGroupUserConfig {
		opt.conds = append(opt.conds, handler.NewCond(cond, value))
		return opt
	}
}

type groupUserProjection struct {
	es handler.EventStore
}

func newGroupUserProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &groupUserProjection{es: config.Eventstore})
}

func (*groupUserProjection) Name() string {
	return GroupUserProjectionTable
}

func (*groupUserProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable(
			groupUserColumns, // handler.NewColumn(GroupUserGroupIDCol, handler.ColumnTypeText),
			handler.NewPrimaryKey(GroupUserInstanceID, GroupUserGroupIDCol, GroupUserUserIDCol),
			handler.WithIndex(handler.NewIndex("user_id", []string{GroupUserUserIDCol})),
			handler.WithIndex(
				handler.NewIndex("gm_instance", []string{GroupUserInstanceID},
					handler.WithInclude(
						GroupUserCreationDate,
						GroupUserChangeDate,
						GroupUserSequence,
						GroupUserResourceOwner,
						GroupUserAttributes,
					),
				),
			),
		),
	)
}

func reduceGroupUserAdded(e groupusers.GroupUserAddedEvent, userResourceOwner string, opts ...reduceGroupUserOpt) (*handler.Statement, error) {
	config := reduceGroupUserConfig{
		cols: []handler.Column{
			handler.NewCol(GroupUserUserIDCol, e.UserID),
			// handler.NewCol(GroupUserGroupIDCol, e.GroupID),
			handler.NewCol(GroupUserUserResourceOwner, userResourceOwner),
			handler.NewCol(GroupUserCreationDate, e.CreatedAt()),
			handler.NewCol(GroupUserChangeDate, e.CreatedAt()),
			handler.NewCol(GroupUserSequence, e.Sequence()),
			handler.NewCol(GroupUserResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(GroupUserInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(GroupUserAttributes, e.Attributes),
		}}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.NewCreateStatement(&e, config.cols), nil
}

func reduceGroupUserChanged(e groupusers.GroupUserChangedEvent, opts ...reduceGroupUserOpt) (*handler.Statement, error) {
	config := reduceGroupUserConfig{
		cols: []handler.Column{
			handler.NewCol(GroupUserChangeDate, e.CreatedAt()),
			handler.NewCol(GroupUserSequence, e.Sequence()),
			handler.NewCol(GroupUserAttributes, e.Attributes),
			handler.NewCol(GroupUserUserIDCol, e.UserID),
		},
		conds: []handler.Condition{
			handler.NewCond(GroupUserInstanceID, e.Aggregate().InstanceID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.NewUpdateStatement(&e, config.cols, config.conds), nil
}

func reduceGroupUserCascadeRemoved(e groupusers.GroupUserCascadeRemovedEvent, opts ...reduceGroupUserOpt) (*handler.Statement, error) {
	config := reduceGroupUserConfig{
		conds: []handler.Condition{
			handler.NewCond(GroupUserInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupUserUserIDCol, e.UserID),
		}}

	for _, opt := range opts {
		config = opt(config)
	}
	return handler.NewDeleteStatement(&e, config.conds), nil
}

func reduceGroupUserRemoved(e eventstore.Event, opts ...reduceGroupUserOpt) (*handler.Statement, error) {
	config := reduceGroupUserConfig{
		conds: []handler.Condition{
			handler.NewCond(GroupUserInstanceID, e.Aggregate().InstanceID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}
	return handler.NewDeleteStatement(e, config.conds), nil
}

func multiReduceGroupUserOwnerRemoved(e eventstore.Event, opts ...reduceGroupUserOpt) func(eventstore.Event) handler.Exec {
	config := reduceGroupUserConfig{
		conds: []handler.Condition{
			handler.NewCond(GroupUserInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupUserResourceOwner, e.Aggregate().ID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}

	return handler.AddDeleteStatement(
		config.conds,
	)
}

func groupUserUserOwnerRemovedConds(e eventstore.Event, opts ...reduceGroupUserOpt) []handler.Condition {
	config := reduceGroupUserConfig{
		conds: []handler.Condition{
			handler.NewCond(GroupUserInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(GroupUserUserResourceOwner, e.Aggregate().ID),
		},
	}

	for _, opt := range opts {
		config = opt(config)
	}
	return config.conds
}

func reduceGroupUserUserOwnerRemoved(e eventstore.Event, opts ...reduceGroupUserOpt) (*handler.Statement, error) {
	return handler.NewDeleteStatement(
		e,
		groupUserUserOwnerRemovedConds(e, opts...),
	), nil
}

func multiReduceGroupUserUserOwnerRemoved(e eventstore.Event, opts ...reduceGroupUserOpt) func(eventstore.Event) handler.Exec {
	return handler.AddDeleteStatement(
		groupUserUserOwnerRemovedConds(e, opts...),
	)
}

func setGroupUserContext(aggregate *eventstore.Aggregate) context.Context {
	return authz.WithInstanceID(context.Background(), aggregate.InstanceID)
}

func (g *groupUserProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  groupusers.AddedEventType,
					Reduce: g.reduceAdded,
				},
				{
					Event:  groupusers.ChangedEventType,
					Reduce: g.reduceChanged,
				},
				{
					Event:  groupusers.CascadeRemovedEventType,
					Reduce: g.reduceCascadeRemoved,
				},
				{
					Event:  groupusers.RemovedEventType,
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
					Reduce: reduceInstanceRemovedHelper(GroupUserInstanceID),
				},
			},
		},
	}
}

func (g *groupUserProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*groupusers.GroupUserAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-chy6O", "reduce.wrong.event.type %s", groupusers.AddedEventType)
	}
	ctx := setGroupUserContext(e.Aggregate())
	userOwner, err := getUserResourceOwner(ctx, g.es, e.Aggregate().InstanceID, e.UserID)
	if err != nil {
		return nil, err
	}
	return reduceGroupUserAdded(
		*e,
		userOwner,
		withGroupUserCol(GroupUserGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupUserProjection) reduceChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*groupusers.GroupUserChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-10XK2", "reduce.wrong.event.type %s", groupusers.ChangedEventType)
	}
	return reduceGroupUserChanged(
		*e,
		withGroupUserCond(GroupUserGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupUserProjection) reduceCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*groupusers.GroupUserCascadeRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-bHe54", "reduce.wrong.event.type %s", groupusers.CascadeRemovedEventType)
	}
	return reduceGroupUserCascadeRemoved(
		*e,
		withGroupUserCond(GroupUserGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupUserProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*groupusers.GroupUserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-fKAOi", "reduce.wrong.event.type %s", groupusers.RemovedEventType)
	}
	return reduceGroupUserRemoved(
		e,
		withGroupUserCond(GroupUserUserIDCol, e.UserID),
		withGroupUserCond(GroupUserGroupIDCol, e.Aggregate().ID),
	)
}

func (g *groupUserProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-aYA60", "reduce.wrong.event.type %s", user.UserRemovedType)
	}
	return reduceGroupUserRemoved(e,
		withGroupUserCond(GroupUserUserIDCol, e.Aggregate().ID),
	)
}

func (g *groupUserProjection) reduceOrgRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-OHVFM", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}
	return handler.NewMultiStatement(
		e,
		multiReduceGroupUserOwnerRemoved(e),
		multiReduceGroupUserUserOwnerRemoved(e),
	), nil
}

func (g *groupUserProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*group.GroupRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-OHVFM", "reduce.wrong.event.type %s", group.GroupRemovedType)
	}
	return reduceGroupUserRemoved(e, withGroupUserCond(GroupUserGroupIDCol, e.Aggregate().ID))
}
