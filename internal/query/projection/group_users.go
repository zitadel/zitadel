package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/group"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

const (
	GroupUsersProjectionTable = "projections.group_users1"

	GroupUsersColumnGroupID       = "group_id"
	GroupUsersColumnUserID        = "user_id"
	GroupUsersColumnResourceOwner = "resource_owner"
	GroupUsersColumnInstanceID    = "instance_id"
	GroupUsersColumnSequence      = "sequence"
	GroupUsersColumnCreationDate  = "creation_date"
)

type groupUsersProjection struct{}

func (g *groupUsersProjection) Name() string {
	return GroupUsersProjectionTable
}

func newGroupUsersProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(groupUsersProjection))
}

func (*groupUsersProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(GroupUsersColumnGroupID, handler.ColumnTypeText),
			handler.NewColumn(GroupUsersColumnUserID, handler.ColumnTypeText),
			handler.NewColumn(GroupUsersColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(GroupUsersColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(GroupUsersColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(GroupUsersColumnCreationDate, handler.ColumnTypeTimestamp),
		},
			handler.NewPrimaryKey(GroupUsersColumnInstanceID, GroupUsersColumnGroupID, GroupUsersColumnUserID),
			handler.WithIndex(handler.NewIndex("user_id", []string{GroupUsersColumnUserID})),
			handler.WithIndex(handler.NewIndex("group_id", []string{GroupUsersColumnGroupID})),
		),
	)
}

func (g *groupUsersProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  group.GroupUserAddedEventType,
					Reduce: g.reduceGroupUserAdded,
				},
				{
					Event:  group.GroupUserRemovedEventType,
					Reduce: g.reduceGroupUserRemoved,
				},
				{
					Event:  group.GroupRemovedEventType,
					Reduce: g.reduceGroupRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: g.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(GroupColumnInstanceID),
				},
			},
		},
	}
}

func (g *groupUsersProjection) reduceGroupUserAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupUserAddedEvent](event)
	if err != nil {
		return nil, err
	}
	// ON CONFLICT DO UPDATE is belt-and-suspenders alongside the eventstore
	// unique constraint registered by GroupUserAddedEvent.UniqueConstraints();
	// the constraint is the primary defense against duplicate memberships,
	// the upsert keeps the projection self-healing if an old duplicate event
	// is ever re-applied.
	return handler.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(GroupUsersColumnInstanceID, nil),
			handler.NewCol(GroupUsersColumnGroupID, nil),
			handler.NewCol(GroupUsersColumnUserID, nil),
		},
		[]handler.Column{
			handler.NewCol(GroupUsersColumnGroupID, e.Aggregate().ID),
			handler.NewCol(GroupUsersColumnUserID, e.UserID),
			handler.NewCol(GroupUsersColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(GroupUsersColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(GroupUsersColumnSequence, e.Sequence()),
			handler.NewCol(GroupUsersColumnCreationDate, handler.OnlySetValueOnInsert(GroupUsersProjectionTable, e.CreationDate())),
		},
	), nil
}

func (g *groupUsersProjection) reduceGroupUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupUserRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupUsersColumnGroupID, e.Aggregate().ID),
			handler.NewCond(GroupUsersColumnUserID, e.UserID),
			handler.NewCond(GroupUsersColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCond(GroupUsersColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupUsersProjection) reduceGroupRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupUsersColumnGroupID, e.Aggregate().ID),
			handler.NewCond(GroupUsersColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCond(GroupUsersColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (g *groupUsersProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.OrgRemovedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(GroupUsersColumnResourceOwner, e.Aggregate().ID),
			handler.NewCond(GroupUsersColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
