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
					Event:  group.GroupUsersAddedEventType,
					Reduce: g.reduceGroupUsersAdded,
				},
				{
					Event:  group.GroupUsersRemovedEventType,
					Reduce: g.reduceGroupUsersRemoved,
				},
			},
		},
		{
			Aggregate: group.AggregateType,
			EventReducers: []handler.EventReducer{
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

func (g *groupUsersProjection) reduceGroupUsersAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupUsersAddedEvent](event)
	if err != nil {
		return nil, err
	}

	stmts := make([]func(eventstore.Event) handler.Exec, 0, len(e.UserIDs))
	for _, userID := range e.UserIDs {
		stmts = append(stmts, handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(GroupUsersColumnGroupID, e.Aggregate().ID),
				handler.NewCol(GroupUsersColumnUserID, userID),
				handler.NewCol(GroupUsersColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(GroupUsersColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(GroupUsersColumnSequence, e.Sequence()),
				handler.NewCol(GroupUsersColumnCreationDate, e.CreationDate()),
			},
		))
	}
	return handler.NewMultiStatement(e, stmts...), nil
}

func (g *groupUsersProjection) reduceGroupUsersRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*group.GroupUsersRemovedEvent](event)
	if err != nil {
		return nil, err
	}
	stmts := make([]func(eventstore.Event) handler.Exec, 0, len(e.UserIDs))
	for _, userID := range e.UserIDs {
		stmts = append(stmts, handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(GroupUsersColumnGroupID, e.Aggregate().ID),
				handler.NewCond(GroupUsersColumnUserID, userID),
				handler.NewCond(GroupUsersColumnInstanceID, e.Aggregate().InstanceID),
			},
		))
	}
	return handler.NewMultiStatement(e, stmts...), nil
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
